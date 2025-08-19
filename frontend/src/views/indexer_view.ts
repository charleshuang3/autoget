import { html, LitElement } from 'lit';
import { customElement, property } from 'lit/decorators.js';

@customElement('indexer-view')
export class IndexerView extends LitElement {
  @property({ type: String })
  index_id = '';

  render() {
    return html`<div>Indexer View: ${this.index_id}</div>`;
  }
}
