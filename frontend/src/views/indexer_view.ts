import { html, LitElement } from 'lit';
import { customElement, property } from 'lit/decorators.js';
import '../components/navbar.ts';

@customElement('indexer-view')
export class IndexerView extends LitElement {
  @property({ type: String })
  index_id = '';

  render() {
    return html`
      <app-navbar .activePage=${this.index_id}></app-navbar>
      <div>Indexer View: ${this.index_id}</div>
    `;
  }
}
