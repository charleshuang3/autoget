import { html, LitElement } from 'lit';
import { customElement, property } from 'lit/decorators.js';
import { consume } from '@lit/context';

import { indexerIdContext } from '../context.ts';
import '../components/navbar.ts';

@customElement('indexer-view')
export class IndexerView extends LitElement {
  @consume({ context: indexerIdContext, subscribe: true })
  @property({ attribute: false })
  public indexer_id = '';

  render() {
    return html`
      <app-navbar .activePage=${this.indexer_id}></app-navbar>
      <div>Indexer View: ${this.indexer_id}</div>
    `;
  }
}
