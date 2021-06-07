package gee

import "testing"

func TestInsert(t *testing.T) {
	tree := NewTrie()
	tree.Insert("/hello/world")
	tree.Insert("/hello")
	tree.Insert("/")
	tree.Insert("/hello/docker")
	tree.Insert("/static/*filepath")

	t.Log(tree.root.part)
	t.Log(tree.root.children[0].part)
	t.Log(tree.root.children[0].children[0].part)
}

func TestInsertRepeatedly(t *testing.T) {
	tree := NewTrie()
	tree.Insert("/hello/world")
	tree.Insert("/:hello/world")

	t.Log(tree.root.part)
	t.Log(tree.root.children[0].part)
	t.Log(tree.root.children[0].children[0].part)
}

func TestAnyUrl(t *testing.T) {
	tree := NewTrie()
	tree.Insert("/p/:lang/doc")

	t.Log(tree.root.part)
	t.Log(tree.root.children[0].part)
	t.Log(tree.root.children[0].children[0].part)
}

func TestSearch(t *testing.T) {
	tree := NewTrie()
	tree.Insert("/hello")
	tree.Insert("/")
	tree.Insert("/static/*filepath")

	searchNode1, params1 := tree.Search("/")
	t.Log(searchNode1.pattern)
	t.Log(params1)

	//searchNode2, params2 := tree.Search("/test/world")
	//t.Log(searchNode2.pattern)
	//t.Log(params2)

	searchNode3, params3 := tree.Search("/static/css/style.css")
	t.Log(searchNode3.pattern)
	t.Log(params3)
}

func TestSearchParam(t *testing.T) {
	tree := NewTrie()
	tree.Insert("/:hello/world")
	tree.Insert("/hello")
	tree.Insert("/")

	searchNode1, params1 := tree.Search("/")
	t.Log(searchNode1.pattern)
	t.Log(params1)

	searchNode2, params2 := tree.Search("/hello")
	t.Log(searchNode2.pattern)
	t.Log(params2)

	searchNode3, params3 := tree.Search("/hello/world")
	t.Log(searchNode3.pattern)
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
