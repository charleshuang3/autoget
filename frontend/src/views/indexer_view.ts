import { html, LitElement, unsafeCSS, css, type TemplateResult, type PropertyValues } from 'lit';
import { customElement, property, state } from 'lit/decorators.js';

import { type Category, fetchIndexerCategories, fetchIndexerResources, type ResourcesResponse } from '../utils/api';
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

  @state()
  private resources: ResourcesResponse | null = null;

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
    await this.fetchIndexerResources();
  }

  protected update(changedProperties: PropertyValues): void {
    if (changedProperties.has('indexerId')) {
      this.fetchIndexerCategories();
      this.fetchIndexerResources();
    }
    if (changedProperties.has('category')) {
      this.fetchIndexerResources();
    }
    super.update(changedProperties);
  }

  private async fetchIndexerCategories() {
    this.categories = await fetchIndexerCategories(this.indexerId);
  }

  private async fetchIndexerResources() {
    if (this.indexerId && this.category) {
      this.resources = await fetchIndexerResources(this.indexerId, this.category, 1);
    } else {
      this.resources = null;
    }
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

          <div class="flex-10 p-4 overflow-y-auto" id="content">
            Indexer View: ${this.indexerId}
            <p>Category: ${this.category}</p>
            ${this.resources
              ? html`<p>Total Resources: ${this.resources.pagination.total}</p>`
              : html`<p>No resources found or loading...</p>`}
          </div>
        </div>
      </div>
    `;
  }
}
