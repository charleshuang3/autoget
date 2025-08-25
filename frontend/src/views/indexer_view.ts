import { html, LitElement, unsafeCSS, type TemplateResult } from 'lit';
import { customElement, property } from 'lit/decorators.js';
import { consume } from '@lit/context';
import { until } from 'lit/directives/until.js';

import { indexerIdContext, indexerDetailsContext, type IndexerDetails, type Category } from '../context.ts';
import '../components/navbar.ts';
import globalStyles from '/src/index.css?inline';

@customElement('indexer-view')
export class IndexerView extends LitElement {
  static styles = [unsafeCSS(globalStyles)];

  @consume({ context: indexerIdContext, subscribe: true })
  @property({ attribute: false })
  public indexer_id = '';

  @consume({ context: indexerDetailsContext, subscribe: true })
  @property({ attribute: false })
  public indexerDetails!: IndexerDetails;

  private renderCategory(category: Category): TemplateResult {
    if (category.subCategories && category.subCategories.length > 0) {
      return html`
        <li>
          <details>
            <summary>${category.name}</summary>
            <ul>
              ${category.subCategories.map((child) => this.renderCategory(child))}
            </ul>
          </details>
        </li>
      `;
    } else {
      return html` <li><a>${category.name}</a></li> `;
    }
  }

  render() {
    return html`
      <div class="flex flex-col h-screen">
        <app-navbar .activePage=${this.indexer_id}></app-navbar>

        <div class="flex flex-row flex-grow">
          <div class="flex-2 bg-base-200" id="left-panel-categories">
            <ul class="menu bg-base-200 rounded-box w-full">
              ${until(
                this.indexerDetails.categories(this.indexer_id).then((categories) => {
                  return categories.map((category) => this.renderCategory(category));
                }),
              )}
            </ul>
          </div>

          <div class="flex-10 p-4" id="content">Indexer View: ${this.indexer_id}</div>
        </div>
      </div>
    `;
  }
}
