import { Router } from '@lit-labs/router';
import { LitElement, html } from 'lit';
import { customElement, property } from 'lit/decorators.js';
import { consume, provide } from '@lit/context';
import { indexersContext, indexIdContext } from './context.ts';

import './views/search_view';
import './views/indexer_view';

@customElement('app-router')
export class AppRouter extends LitElement {
  @consume({ context: indexersContext, subscribe: true })
  @property({ attribute: false })
  public indexers: string[] = [];

  @provide({ context: indexIdContext })
  @property({ attribute: false })
  public index_id = '';

  private router = new Router(this, [
    {
      path: '/',
      render: () => html`<div>Loading...</div>`,
      enter: async () => {
        let newUrl = `/indexers/${this.indexers[0]}`;
        this.router.goto(newUrl);
        history.replaceState(null, '', newUrl);
        return false;
      },
    },
    { path: '/search', render: () => html`<search-view></search-view>` },
    {
      path: '/indexers/:id',
      render: ({ id }) => {
        this.index_id = id || '';
        return html`<indexer-view></indexer-view>`;
      },
    },
  ]);

  render() {
    return this.router.outlet();
  }
}
