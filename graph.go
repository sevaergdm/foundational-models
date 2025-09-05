package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"oss.terrastruct.com/d2/d2format"
	"oss.terrastruct.com/d2/d2graph"
	"oss.terrastruct.com/d2/d2layouts/d2dagrelayout"
	"oss.terrastruct.com/d2/d2lib"
	"oss.terrastruct.com/d2/d2oracle"
	"oss.terrastruct.com/d2/d2renderers/d2svg"
	dlog "oss.terrastruct.com/d2/lib/log"
	"oss.terrastruct.com/d2/lib/textmeasure"
)


func createNode(entity FoundationalModel, g *d2graph.Graph) (*d2graph.Graph, error) {
	updatedGraph := g

	newG, newKey, err := d2oracle.Create(updatedGraph, nil, entity.Name)
	if err != nil {
		return nil, err
	}
	updatedGraph = newG

	shape := "sql_table"
	newG, err = d2oracle.Set(updatedGraph, nil, fmt.Sprintf("%s.shape", newKey), nil, &shape)
	if err != nil {
		return nil, err
	}
	updatedGraph = newG

	for _, attr := range entity.Attributes {
		newG, err = d2oracle.Set(updatedGraph, nil, fmt.Sprintf("%s.%s", entity.Name, attr.Name), nil, &attr.DataType)
		if err != nil {
			return nil, err
		}
		updatedGraph = newG
	}

	return updatedGraph, nil	
}

func createEdges(entity FoundationalModel, g *d2graph.Graph) (*d2graph.Graph, error) {
	updatedGraph := g

	for _, rel := range entity.Relationships {
		var err error
		source := fmt.Sprintf("%s.%s", entity.Name, rel.RelatedAttribute)
		target := fmt.Sprintf("%s.%s", rel.RelatedEntity, rel.RelatedAttribute)
		updatedGraph, _, err = d2oracle.Create(updatedGraph, nil, fmt.Sprintf("(%s -> %s)", source, target))
		if err != nil {
			return nil, err
		}
	}
	return updatedGraph, nil
}

func (cfg *apiConfig) graphBuilder() error {
	ctx := dlog.WithDefault(context.Background())

	_, graph, err := d2lib.Compile(ctx, "", nil, nil)
	if err != nil {
		return err
	}


	for _, v := range cfg.entitiesCache {
		graph, err = createNode(v, graph)
		if err != nil {
			return fmt.Errorf("failed to create node for entity '%s': %w", v.Name, err)
		}
	}

	for _, v := range cfg.entitiesCache {
		graph, err = createEdges(v, graph)
		if err != nil {
			return fmt.Errorf("failed to create edges for entity '%s': %w", v.Name, err)
		}
	}

	ruler, _ := textmeasure.NewRuler()
	script := d2format.Format(graph.AST)

	layoutResolver := func(engine string) (d2graph.LayoutGraph, error) {
		return d2dagrelayout.DefaultLayout, nil
	}

	diagram, _, err := d2lib.Compile(ctx, script, &d2lib.CompileOptions{
		LayoutResolver: layoutResolver,
		Ruler: ruler,
	}, nil)
	if err != nil {
		return err
	}

	padding := int64(d2svg.DEFAULT_PADDING)
	out, _ := d2svg.Render(diagram, &d2svg.RenderOpts{
		Pad: &padding,
	})

	svgPath := filepath.Join("svgs", "final.svg")
	err = os.WriteFile(svgPath, out, 0600)
	if err != nil {
		return err
	}

	return os.WriteFile("out.d2", []byte(script), 0600)
}
