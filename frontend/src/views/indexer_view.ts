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

  @property({ type: Number })
  public currentPage: number = 1;

  @state()
  private totalPages: number = 1;

  private formatBytes(bytes: number, decimals: number = 2): string {
    if (bytes === 0) return '0 Bytes';

    const k = 1024;
    const dm = decimals < 0 ? 0 : decimals;
    const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];

    const i = Math.floor(Math.log(bytes) / Math.log(k));

    return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i];
  }

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
      this.resources = null;
      this.currentPage = 1; // Reset page when indexerId changes
      await this.fetchIndexerCategories();
      await this.fetchIndexerResources();
    }
    if (changedProperties.has('category')) {
      this.resources = null;
      this.currentPage = 1; // Reset page when category changes
      this.fetchIndexerResources();
    }
    if (changedProperties.has('currentPage')) {
      this.resources = null;
      this.fetchIndexerResources();
    }
    super.update(changedProperties);
  }

  private async fetchIndexerCategories() {
    this.categories = await fetchIndexerCategories(this.indexerId);
  }

  private async fetchIndexerResources() {
    if (this.indexerId && this.category) {
      const response = await fetchIndexerResources(this.indexerId, this.category, this.currentPage);
      if (response) {
        this.resources = response;
        this.totalPages = response.pagination.totalPages;
      } else {
        this.resources = null;
        this.totalPages = 1;
      }
    } else {
      this.resources = null;
      this.totalPages = 1;
    }
  }

  private handlePageChange(page: number) {
    if (page >= 1 && page <= this.totalPages) {
      this.currentPage = page;
    }
  }

  private renderPagination(): TemplateResult | null {
    if (this.totalPages <= 1) {
      return null;
    }

    const pages: (number | string)[] = [];
    const maxPagesToShow = 5;
    const half = Math.floor(maxPagesToShow / 2);

    let startPage = Math.max(1, this.currentPage - half);
    let endPage = Math.min(this.totalPages, this.currentPage + half);

    if (endPage - startPage + 1 < maxPagesToShow) {
      if (this.currentPage <= half) {
        endPage = Math.min(this.totalPages, maxPagesToShow);
      } else if (this.currentPage + half >= this.totalPages) {
        startPage = Math.max(1, this.totalPages - maxPagesToShow + 1);
      }
    }

    if (startPage > 1) {
      pages.push('<');
    }

    for (let i = startPage; i <= endPage; i++) {
      pages.push(i);
    }

    if (endPage < this.totalPages) {
      pages.push('>');
    }

    return html`
      <div class="flex justify-center my-4">
        <div class="join">
          ${pages.map((page) => {
            const isActive = page === this.currentPage;
            const isDisabled =
              (page === '<' && this.currentPage === 1) || (page === '>' && this.currentPage === this.totalPages);
            const buttonClass = `join-item btn ${isActive ? 'btn-active' : ''} ${isDisabled ? 'btn-disabled' : ''}`;

            if (typeof page === 'number') {
              return html`<button class="${buttonClass}" @click=${() => this.handlePageChange(page)}>${page}</button>`;
            } else if (page === '<') {
              return html`<button class="${buttonClass}" @click=${() => this.handlePageChange(this.currentPage - 1)}>
                &laquo;
              </button>`;
            } else if (page === '>') {
              return html`<button class="${buttonClass}" @click=${() => this.handlePageChange(this.currentPage + 1)}>
                &raquo;
              </button>`;
            }
            return null;
          })}
        </div>
      </div>
    `;
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
            ${this.renderPagination()}
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
                        <div class="p-2">
                          <h3 class="text font-medium line-clamp-4 text-balance break-all">
                            ${resource.title || 'Untitled Resource'}
                          </h3>
                          ${resource.title2
                            ? html`<p class="text font-normal line-clamp-4 text-balance break-all">
                                ${resource.title2}
                              </p>`
                            : ''}
                          <div class="flex flex-wrap gap-1 mt-2">
                            <span class="badge badge-outline badge-primary">${resource.category}</span>
                            <span class="badge badge-outline badge-secondary">${this.formatBytes(resource.size)}</span>
                            ${resource.resolution
                              ? html`<span class="badge badge-outline badge-info">${resource.resolution}</span>`
                              : ''}
                          </div>
                        </div>
                      </div>
                    `,
                  )
                : html`<p>No resources found or loading...</p>`}
            </div>
            ${this.renderPagination()}
          </div>
        </div>
      </div>
    `;
  }
}
