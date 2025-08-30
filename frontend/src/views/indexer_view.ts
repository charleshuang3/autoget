import { html, LitElement, unsafeCSS, css, type TemplateResult } from 'lit';
import { customElement, property, state } from 'lit/decorators.js';

import { type Category, fetchIndexerCategories } from '../utils/api';
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

  @property({ type: String })
  public indexerId: string = '';

  @state()
  private categories: Category[] = [];

  @property({ type: String })
  public category: string = '';

  private renderCategory(category: Category): TemplateResult {
    const isActive = this.category === category.id;
    const activeClass = isActive ? 'menu-active' : '';

    if (category.subCategories && category.subCategories.length > 0) {
      return html`
        <li>
          <a class="${activeClass}" href="/indexers/${this.indexerId}/${category.id}">${category.name}</a>
          <ul>
            ${category.subCategories.map((child) => this.renderCategory(child))}
          </ul>
        </li>
      `;
    } else {
      return html`<li>
        <a class="${activeClass}" href="/indexers/${this.indexerId}/${category.id}">${category.name}</a>
      </li> `;
    }
  }

  async connectedCallback() {
    super.connectedCallback();
    await this.fetchIndexerCategories();
    await this.setDefaultCategory();
  }

  private async setDefaultCategory() {
    if (!this.category) {
      this.category = this.categories[0].id;
    }
  }
  private async fetchIndexerCategories() {
    this.categories = await fetchIndexerCategories(this.indexerId);
  }

  render() {
    return html`
      <div class="flex flex-col h-screen">
        <app-navbar .activePage=${this.indexerId}></app-navbar>

        <div class="flex flex-row flex-grow overflow-hidden">
          <div class="flex-2 bg-base-200 overflow-y-auto" id="left-panel-categories">
            <ul class="menu bg-base-200 rounded-box w-full">
              ${this.categories.map((category) => this.renderCategory(category))}
            </ul>
          </div>

          <div class="flex-10 p-4 overflow-y-auto" id="content">Indexer View: ${this.indexerId}</div>
        </div>
      </div>
    `;
  }
}
