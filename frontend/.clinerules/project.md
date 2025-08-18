## Styling and Component Structure

- Use **Lit** framework (`@lit/reactive-element` or `lit`) for all web components.
- Write code in **TypeScript**, with strict typings enabled.
- Use **DaisyUI** classes for design and styling: e.g., `btn`, `card`, etc.
- Implement layout using **CSS Flexbox** only—prefer `flex`, `flex-col`, `flex-row`, `flex-grow`, `items-`, `justify-` identifiers.
- Use `@vaadin/router` for Single Page Application (SPA) routing.

## Code Conventions

- Structure each component:
  1. Props & types at top using `interface`.
  2. component class definition
  3. `render()` method
  4. CSS styles at the bottom inside `static styles = css\`...\`;`

- Follow DaisyUI’s semantic class names (e.g., `btn-primary`, `card-body`), and avoid custom CSS unless absolutely necessary.

- Keep component files focused, naming them clearly: `MyComponent.ts`.

- Prioritize readability and consistency over brevity.
