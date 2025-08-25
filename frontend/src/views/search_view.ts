import { html, LitElement } from 'lit';
import { customElement } from 'lit/decorators.js';

import '../components/navbar.ts';

@customElement('search-view')
export class SearchView extends LitElement {
  render() {
    return html`
      <app-navbar activePage="search"></app-navbar>
      <div>Search View</div>
    `;
  }
}
