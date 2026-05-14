// Package template provides placeholder expansion for envchain variable values.
//
// Variable values may reference other variables within the same chain using
// the ${VAR_NAME} syntax. When a chain is rendered, all resolvable placeholders
// are expanded in-place. Unresolvable references (variables not defined in the
// chain) are left unchanged so that downstream tools (e.g. shells) can handle
// them.
//
// Circular references are detected eagerly and returned as errors to prevent
// infinite recursion.
//
// Example:
//
//	renderer := template.NewRenderer(store)
//	vars, err := renderer.Render("myproject")
//	if err != nil {
//		log.Fatal(err)
//	}
//	for k, v := range vars {
//		fmt.Printf("%s=%s\n", k, v)
//	}
package template
