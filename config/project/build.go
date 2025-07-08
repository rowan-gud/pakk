package project

import (
	"fmt"
	"log/slog"
	"slices"
	"sync"

	"github.com/rowan-gud/pakk/collections"
	"github.com/rowan-gud/pakk/config/mod"
)

func (p *Project) Build() error {
	buildPlan := make(map[string][][]*mod.Mod)

	for key, node := range p.modules.Nodes() {
		if node.Data().Bin != nil {
			if !p.checkAcyclicRoot(node, nil) {
				return fmt.Errorf("cannot build module %s: cyclic graph", key)
			}

			plan, err := p.planBuildRoot(node)
			if err != nil {
				return fmt.Errorf("failed to construct build plan: %w", err)
			}
			buildPlan[key] = plan
		}
	}

	for key, layers := range buildPlan {
		if len(layers) == 0 {
			continue
		}

		for _, layer := range layers {
			mut := sync.Mutex{}
			errChan := make(chan error)

			for _, mod := range layer {
				p.logger.Info("building module",
					slog.String("module", mod.Path()),
				)

				go func(key string, errChan chan error) {
					err := mod.Build()
					if err != nil {
						errChan <- err
					}

					sum, err := mod.Sum()
					if err != nil {
						errChan <- err
					}

					mut.Lock()
					p.lockFile.Modules[mod.Path()] = sum
					mut.Unlock()

					errChan <- nil
				}(key, errChan)
			}

			for range len(layer) {
				err := <-errChan
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (p *Project) checkAcyclicRoot(node *collections.DirectedGraphNode[string, *mod.Mod], seen *collections.Set[*collections.DirectedGraphNode[string, *mod.Mod]]) bool {
	if node == nil {
		return true
	}

	if seen == nil {
		seen = collections.NewSet[*collections.DirectedGraphNode[string, *mod.Mod]]()
	}

	for _, n := range node.Edges() {
		if seen.Has(n) {
			return false
		}

		seen.Add(n)
		isAcyclic := p.checkAcyclicRoot(n, seen)
		if !isAcyclic {
			return false
		}
		seen.Delete(n)
	}

	return true
}

func (p *Project) planBuildRoot(root *collections.DirectedGraphNode[string, *mod.Mod]) ([][]*mod.Mod, error) {
	buildTree, err := p.planBuildNode(root)
	if err != nil {
		return nil, err
	}

	if buildTree == nil {
		return nil, nil
	}

	numLayers := buildTree.Height()
	layers := make([][]*mod.Mod, numLayers)

	p.planBuildLayers(buildTree, 0, layers)

	seen := collections.NewSet[*mod.Mod]()
	deduped := make([][]*mod.Mod, 0, numLayers)

	for _, layer := range slices.Backward(layers) {
		newLayer := make([]*mod.Mod, 0, len(layer))

		for _, mod := range layer {
			if seen.Has(mod) {
				continue
			}

			newLayer = append(newLayer, mod)
			seen.Add(mod)
		}

		if len(newLayer) > 0 {
			deduped = append(deduped, newLayer)
		}
	}

	return deduped, nil
}

func (p *Project) planBuildLayers(node *collections.Tree[*mod.Mod], layerIdx int, layers [][]*mod.Mod) {
	layers[layerIdx] = append(layers[layerIdx], node.Data())

	for _, n := range node.Children() {
		p.planBuildLayers(n, layerIdx+1, layers)
	}
}

func (p *Project) planBuildNode(node *collections.DirectedGraphNode[string, *mod.Mod]) (*collections.Tree[*mod.Mod], error) {
	deps := node.Edges()

	toBuild := make([]*collections.Tree[*mod.Mod], 0, len(deps))

	for _, dep := range deps {
		mods, err := p.planBuildNode(dep)
		if err != nil {
			return nil, err
		}

		if mods != nil {
			toBuild = append(toBuild, mods)
		}
	}

	mod := node.Data()

	sum, err := mod.Sum()
	if err != nil {
		return nil, fmt.Errorf("failed to compute checksum for module %s: %w", node.Key, err)
	}

	if len(toBuild) > 0 || sum != p.lockFile.Modules[node.Key] {
		return collections.NewTree(mod, toBuild...), nil
	}

	return nil, nil
}
