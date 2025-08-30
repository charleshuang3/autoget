import { html, LitElement, unsafeCSS, css, type TemplateResult } from 'lit';
import { customElement, property } from 'lit/decorators.js';
import { consume } from '@lit/context';
import { until } from 'lit/directives/until.js';

import { indexerIdContext, indexerDetailsContext, type IndexerDetails, type Category } from '../context.ts';
import '../components/navbar.ts';
import globalStyles from '/src/index.css?inline';

@customElement('indexer-view')
export class IndexerView extends LitElement {
  static styles = [
    unsafeCSS(globalStyles),
    css`
      #left-panel-categories li ul {
        margin-inline-start: 0.75rem;
        padding-left: 0;
      }
      #left-panel-categories li a {
        padding-left: 0.5rem;
      }
    `,
  ];

  @consume({ context: indexerIdContext, subscribe: true })
  @property({ attribute: false })
  public indexer_id = '';

  @consume({ context: indexerDetailsContext, subscribe: true })
  @property({ attribute: false })
  public indexerDetails!: IndexerDetails;

  @property({ type: String })
  public category: string = '';

  private async setDefaultCategory() {
    if (!this.category) {
      const categories = await this.indexerDetails.categories(this.indexer_id);
      if (categories.length > 0) {
        this.category = categories[0].id;
      }
    }
  }

  private renderCategory(category: Category): TemplateResult {
    const isActive = this.category === category.id;
    const activeClass = isActive ? 'menu-active' : '';

    if (category.subCategories && category.subCategories.length > 0) {
      return html`
        <li>
          <a class="${activeClass}" href="/indexers/${this.indexer_id}/${category.id}">${category.name}</a>
          <ul>
            ${category.subCategories.map((child) => this.renderCategory(child))}
          </ul>
        </li>
      `;
    } else {
      return html`<li><a class="${activeClass}" href="/indexers/${this.indexer_id}/${category.id}">${category.name}</a></li> `;
    }
  }

  async connectedCallback() {
    super.connectedCallback();
    await this.setDefaultCategory();
  }

  render() {
    return html`
      <div class="flex flex-col h-screen">
        <app-navbar .activePage=${this.indexer_id}></app-navbar>

        <div class="flex flex-row flex-grow overflow-hidden">
          <div class="flex-2 bg-base-200 overflow-y-auto" id="left-panel-categories">
            <ul class="menu bg-base-200 rounded-box w-full">
              ${until(
                this.indexerDetails.categories(this.indexer_id).then((categories) => {
                  return categories.map((category) => this.renderCategory(category));
                }),
              )}
            </ul>
          </div>

          <div class="flex-10 p-4 overflow-y-auto" id="content">Indexer View: ${this.indexer_id}</div>
        </div>
      </div>
    `;
  }
}
