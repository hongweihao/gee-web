package gee

import "testing"

func TestInsert(t *testing.T) {
	tree := NewTrie()
	tree.Insert("/hello/world")
	tree.Insert("/hello")
	tree.Insert("/")
	tree.Insert("/hello/docker")
	tree.Insert("/static/*filepath")

	t.Log(tree.Root.Part)
	t.Log(tree.Root.Children[0].Part)
	t.Log(tree.Root.Children[0].Children[0].Part)
}

func TestInsertRepeatedly(t *testing.T) {
	tree := NewTrie()
	tree.Insert("/hello/world")
	tree.Insert("/:hello/world")

	t.Log(tree.Root.Part)
	t.Log(tree.Root.Children[0].Part)
	t.Log(tree.Root.Children[0].Children[0].Part)
}

func TestAnyUrl(t *testing.T) {
	tree := NewTrie()
	tree.Insert("/p/:lang/doc")

	t.Log(tree.Root.Part)
	t.Log(tree.Root.Children[0].Part)
	t.Log(tree.Root.Children[0].Children[0].Part)
}

func TestSearch(t *testing.T) {
	tree := NewTrie()
	tree.Insert("/hello")
	tree.Insert("/")
	tree.Insert("/static/*filepath")

	searchNode1, params1 := tree.Search("/")
	t.Log(searchNode1.Pattern)
	t.Log(params1)

	//searchNode2, params2 := tree.Search("/test/world")
	//t.Log(searchNode2.Pattern)
	//t.Log(params2)

	searchNode3, params3 := tree.Search("/static/css/style.css")
	t.Log(searchNode3.Pattern)
	t.Log(params3)
}

func TestSearchParam(t *testing.T) {
	tree := NewTrie()
	tree.Insert("/:hello/world")
	tree.Insert("/hello")
	tree.Insert("/")

	searchNode1, params1 := tree.Search("/")
	t.Log(searchNode1.Pattern)
	t.Log(params1)

	searchNode2, params2 := tree.Search("/hello")
	t.Log(searchNode2.Pattern)
	t.Log(params2)

	searchNode3, params3 := tree.Search("/hello/world")
	t.Log(searchNode3.Pattern)
	t.Log(params3)
}

func TestParsePattern(t *testing.T) {
	pattern := "/"
	tree := NewTrie()
	parts := tree.parsePattern(pattern)
	t.Log(len(parts))
}

func TestGetParams(t *testing.T) {
	tree := NewTrie()
	params := tree.getParams("/:lang/:hello/test", "/go/hello")
	t.Log(params)

	params1 := tree.getParams("/static/*filepath", "/static/css/style.css")
	t.Log(params1)
}
