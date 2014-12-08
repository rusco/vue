package main

import (
	"github.com/gopherjs/gopherjs/js"
)

const (
	STORAGE_KEY = "todos-vuejs"
	TITLE       = "title"
	COMPLETED   = "completed"

	TODOS      = "todos"
	NEWTODO    = "newTodo"
	EDITEDTODO = "editedTodo"
)

func Save(val interface{}) {
	jsonStr := js.Global.Get("JSON").Call("stringify", val)
	js.Global.Get("localStorage").Call("setItem", STORAGE_KEY, jsonStr)
}
func Fetch() js.Object {

	item := js.Global.Get("localStorage").Call("getItem", STORAGE_KEY)
	if item.IsNull() {
		return nil
	}
	obj := js.Global.Get("JSON").Call("parse", item)
	return obj
}

func log(val ...interface{}) {
	js.Global.Get("console").Call("log", val...)
}

func init() {
	if Fetch() == nil {
		Save(js.S{js.M{TITLE: "Program in Gopher.js", COMPLETED: false}, js.M{TITLE: "Finish another TodoMvc Sample", COMPLETED: true}})
	}
}

func main() {

	app := js.M{
		"el": "#todoapp",
		"data": js.M{
			TODOS:          Fetch(),
			NEWTODO:        "",
			EDITEDTODO:     nil,
			"activeFilter": "all",
			"filters": js.M{
				"all": func() interface{} {
					return true
				},
				"active": func(todo js.Object) interface{} {
					return !todo.Get(COMPLETED).Bool()
				},
				COMPLETED: func(todo js.Object) interface{} {
					return todo.Get(COMPLETED).Bool()
				},
			},
		},
		"ready": func() {

			js.This.Call("$watch", TODOS, func(newValue js.S) {
				Save(newValue)
			}, true)

		},
		"directives": js.M{
			"todo-focus": func(v js.Object) {

				if v.Bool() == false {
					return
				}

				el := js.This.Get("el")
				js.Global.Call("setTimeout", func() {
					el.Call("focus")
				}, 0)
			},
		},
		"filters": js.M{
			"filterTodos": func(todos js.Object) js.Object {

				filters := js.This.Get("filters")
				activeFilter := js.This.Get("activeFilter").Str()
				filterTodos := todos.Call("filter", filters.Get(activeFilter))

				return filterTodos
			},
		},
		"computed": js.M{
			"remaining": func() js.Object {

				activeFilterFn := js.This.Get("filters").Get("active")
				activeTodosLen := js.This.Get(TODOS).Call("filter", activeFilterFn).Get("length")

				return activeTodosLen
			},
			"allDone": js.M{
				"get": func() interface{} {
					return js.This.Get("remaining").Int() == 0
				},
				"set": func(done js.Object) {
					js.This.Get(TODOS).Call("forEach", func(item js.Object) {
						item.Set(COMPLETED, done.Bool())
					})
				},
			},
		},
		"methods": js.M{
			"addTodo": func() {

				newToDoStr := js.This.Get(NEWTODO).Call("trim").Str()

				if len(newToDoStr) == 0 {
					return
				}
				js.This.Get(TODOS).Call("push", js.M{TITLE: newToDoStr, COMPLETED: false})
				js.This.Set(NEWTODO, "")
			},
			"removeTodo": func(todo js.Object) {
				js.This.Get(TODOS).Call("$remove", todo.Get("$data"))
			},
			"editTodo": func(todo js.Object) {

				js.This.Set("beforeEditCache", todo.Get(TITLE))
				js.This.Set(EDITEDTODO, todo)
			},
			"doneEdit": func(todo js.Object) interface{} {

				if len(js.This.Get(EDITEDTODO).Str()) <= 0 {
					return true
				}

				js.This.Set(EDITEDTODO, nil)
				todo.Set(TITLE, todo.Get(TITLE).Call("trim").Str())

				if len(todo.Get(TITLE).Str()) == 0 {
					js.This.Call("removeTodo", todo)
				}
				return true
			},
			"cancelEdit": func(todo js.Object) {

				js.This.Set(EDITEDTODO, nil)
				todo.Set(TITLE, js.This.Get("beforeEditCache"))
			},
			"removeCompleted": func() {

				activeFilterFn := js.This.Get("filters").Get("active")
				activeTodos := js.This.Get(TODOS).Call("filter", activeFilterFn)
				js.This.Set(TODOS, activeTodos)
			},
		},
	}

	vue := js.Global.Get("Vue").New(app)
	js.Global.Set("app", vue)

	router := js.Global.Get("Router").New()
	keys := js.Global.Get("Object").Call("keys", vue.Get("filters"))
	keys.Call("forEach", func(filter js.Object) {
		router.Call("on", filter, func() {
			vue.Set("activeFilter", filter)
		})
	})

	router.Call("configure", js.M{
		"notfound": func() {
			js.Global.Get("location").Set("hash", "")
			vue.Set("activeFilter", "all")
		},
	})
	router.Call("init")
}
