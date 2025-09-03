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

      body::-webkit-scrollbar {
        width: 8px;
      }

      body::-webkit-scrollbar-track {
        background: #1f2937;
      }

      body::-webkit-scrollbar-thumb {
        background-color: #4b5563;
        border-radius: 20px;
        border: 2px solid #1f2937;
      }

      /* To prevent items from splitting across columns */
      .break-inside-avoid-column {
        break-inside: avoid-column;
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

  protected async update(changedProperties: PropertyValues): Promise<void> {
    if (changedProperties.has('indexerId')) {
      await this.fetchIndexerCategories();
      await this.fetchIndexerResources();
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
    this.resources = null;
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
            <header class="text-center mb-12">
              <h1 class="text-5xl md:text-6xl font-extrabold tracking-tight mb-4">${this.indexerId} Resources</h1>
              <p class="text-xl md:text-2xl text-gray-300 max-w-2xl mx-auto">Category: ${this.category}</p>
            </header>

            <div id="masonry-container" class="columns-1 sm:columns-2 md:columns-3 lg:columns-4 xl:columns-5 gap-2">
              ${this.resources && this.resources.resources && this.resources.resources.length > 0
                ? this.resources.resources.map(
                    (resource) => html`
                      <div
                        class="image-card rounded-lg overflow-hidden shadow-lg border border-gray-700 bg-gray-100 break-inside-avoid-column mb-2"
                      >
                        ${resource.images && resource.images.length > 0
                          ? html`<img
                              src="${resource.images[0]}"
                              alt="${resource.title || 'Resource image'}"
                              class="w-full h-auto object-cover rounded-lg"
                              loading="lazy"
                            />`
                          : ''}
                        <div class="p-4">
                          <h3 class="text font-medium">${resource.title || 'Untitled Resource'}</h3>
                        </div>
                      </div>
                    `,
                  )
                : html`<p>No resources found or loading...</p>`}
            </div>
          </div>
        </div>
      </div>
    `;
  }
}
