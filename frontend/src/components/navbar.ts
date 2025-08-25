import { LitElement, html, unsafeCSS } from 'lit';
import { customElement, property } from 'lit/decorators.js';
import { until } from 'lit/directives/until.js';
import { consume } from '@lit/context';

import { indexersContext, indexerDetailsContext, type IndexerDetails, type Category } from '../context.ts';
import globalStyles from '/src/index.css?inline';

@customElement('app-navbar')
export class AppNavbar extends LitElement {
  @consume({ context: indexersContext, subscribe: true })
  @property({ attribute: false })
  public indexers: string[] = [];

  @consume({ context: indexerDetailsContext, subscribe: true })
  @property({ attribute: false })
  public indexerDetails!: IndexerDetails;

  static styles = [unsafeCSS(globalStyles)];

  @property({ type: String })
  activePage = '';

  render() {
    this.indexerDetails.categories(this.activePage).then((categories) => {
      console.log(categories);
    });

    return html`
      <div class="navbar bg-base-200">
        <div class="navbar-start">
          <a href="/" class="btn-ghost">
            <img src="/icon.svg" alt="Icon" class="w-8 h-8" />
          </a>
          <div role="tablist" class="tabs tabs-border">
            ${until(
              this.indexers.map((indexer) => {
                const isActive = this.activePage === indexer;
                return html`<a href="/indexers/${indexer}" class="tab ${isActive ? 'tab-active' : ''}" role="tab"
                  >${indexer}</a
                >`;
              }),
            )}
          </div>
        </div>
        <div class="navbar-end">
          <a href="/search" class="btn btn-ghost ${this.activePage === 'search' ? 'btn-active' : ''}">Search</a>
        </div>
      </div>
    `;
  }
}
