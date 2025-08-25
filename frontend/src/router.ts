import { Router } from '@lit-labs/router';
import { LitElement, html } from 'lit';
import { customElement, property } from 'lit/decorators.js';
import { consume } from '@lit/context';
import { indexersContext } from './context.ts';

import './views/search_view';
import './views/indexer_view';

@customElement('app-router')
export class AppRouter extends LitElement {
  @consume({ context: indexersContext, subscribe: true })
  @property({ attribute: false })
  public indexers: string[] = [];

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
      render: ({ id }) => html`<indexer-view .index_id=${id}></indexer-view>`,
    },
  ]);

  render() {
    return this.router.outlet();
  }
}
