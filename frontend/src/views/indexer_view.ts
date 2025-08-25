import { html, LitElement } from 'lit';
import { customElement, property } from 'lit/decorators.js';
import { consume } from '@lit/context';
import { indexIdContext } from '../context.ts';
import '../components/navbar.ts';

@customElement('indexer-view')
export class IndexerView extends LitElement {
  @consume({ context: indexIdContext, subscribe: true })
  @property({ attribute: false })
  public index_id = '';

  render() {
    return html`
      <app-navbar .activePage=${this.index_id}></app-navbar>
      <div>Indexer View: ${this.index_id}</div>
    `;
  }
}
