import { html, LitElement, unsafeCSS, css, type TemplateResult, type PropertyValues } from 'lit';
import { customElement, property, state } from 'lit/decorators.js';
import { DateTime } from 'luxon';
import 'iconify-icon';

import { fetchIndexerResources, type ResourcesResponse } from '../utils/api';
import globalStyles from '/src/index.css?inline';

@customElement('resource-list')
export class ResourceList extends LitElement {
  static styles = [
    unsafeCSS(globalStyles),
    css`
      /* To prevent items from splitting across columns */
      .break-inside-avoid-column {
        break-inside: avoid-column;
      }
    `,
  ];

  @property({ type: String })
  public indexerId: string = '';

  @property({ type: String })
  public category: string = '';

  @property({ type: Number })
  public page: number = 1;

  @state()
  private resources: ResourcesResponse | null = null;

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

  private formatCreatedDate(timestamp: number): string {
    const createdDate = DateTime.fromSeconds(timestamp, { zone: 'utc' });
    const now = DateTime.now();
    const diff = now.diff(createdDate, ['minutes', 'hours', 'days', 'weeks']).toObject();

    if (diff.weeks && diff.weeks >= 1) {
      return createdDate.toFormat('yyyy-MM-dd');
    } else if (diff.days && diff.days >= 1) {
      return `${Math.floor(diff.days)} day${Math.floor(diff.days) === 1 ? '' : 's'} ago`;
    } else if (diff.hours && diff.hours >= 1) {
      return `${Math.floor(diff.hours)} hour${Math.floor(diff.hours) === 1 ? '' : 's'} ago`;
    } else if (diff.minutes && diff.minutes >= 1) {
      return `${Math.floor(diff.minutes)} min${Math.floor(diff.minutes) === 1 ? '' : 's'} ago`;
    } else {
      return 'just now';
    }
  }

  protected async update(changedProperties: PropertyValues): Promise<void> {
    super.update(changedProperties);

    if (changedProperties.has('indexerId') || changedProperties.has('category') || changedProperties.has('page')) {
      await this.fetchIndexerResources();
    }
  }

  private async fetchIndexerResources() {
    if (this.indexerId && this.category) {
      const response = await fetchIndexerResources(this.indexerId, this.category, this.page);
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
      const url = new URL(window.location.href);
      url.searchParams.set('page', page.toString());
      window.history.pushState({}, '', url.toString());
      window.dispatchEvent(new PopStateEvent('popstate'));
    }
  }

  private renderPagination(): TemplateResult | null {
    if (this.totalPages <= 1) {
      return null;
    }

    const pages: (number | string)[] = [];
    const maxPagesToShow = 5;
    const half = Math.floor(maxPagesToShow / 2);

    let startPage = Math.max(1, this.page - half);
    let endPage = Math.min(this.totalPages, this.page + half);

    if (endPage - startPage + 1 < maxPagesToShow) {
      if (this.page <= half) {
        endPage = Math.min(this.totalPages, maxPagesToShow);
      } else if (this.page + half >= this.totalPages) {
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
            const isActive = page === this.page;
            const isDisabled =
              (page === '<' && this.page === 1) || (page === '>' && this.page === this.totalPages);
            const buttonClass = `join-item btn ${isActive ? 'btn-active' : ''} ${isDisabled ? 'btn-disabled' : ''}`;

            if (typeof page === 'number') {
              return html`<button class="${buttonClass}" @click=${() => this.handlePageChange(page)}>${page}</button>`;
            } else if (page === '<') {
              return html`<button class="${buttonClass}" @click=${() => this.handlePageChange(this.page - 1)}>
                &laquo;
              </button>`;
            } else if (page === '>') {
              return html`<button class="${buttonClass}" @click=${() => this.handlePageChange(this.page + 1)}>
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
                    <h3 class="text font-medium line-clamp-4 text-balance break-all border-b-1 border-b-gray-400">
                      ${resource.title || 'Untitled Resource'}
                    </h3>
                    ${resource.title2
                      ? html`<p
                          class="text font-normal line-clamp-4 text-balance break-all border-b-1 border-b-gray-400"
                        >
                          ${resource.title2}
                        </p>`
                      : ''}
                    <div class="flex flex-wrap gap-1 mt-1 mb-1 pb-1 border-b-1 border-b-gray-400">
                      <span class="badge badge-outline badge-primary">${resource.category}</span>
                      <span class="badge badge-outline badge-secondary">${this.formatBytes(resource.size)}</span>
                      ${resource.resolution
                        ? html`<span class="badge badge-outline badge-info">${resource.resolution}</span>`
                        : ''}
                      ${resource.free ? html`<span class="badge badge-success">Free</span>` : ''}
                      <span
                        class="badge ${DateTime.now().diff(
                          DateTime.fromSeconds(resource.createdDate, { zone: 'utc' }),
                          'weeks',
                        ).weeks < 1
                          ? 'badge-accent'
                          : 'badge-neutral'}"
                      >
                        <iconify-icon icon="mingcute:time-line"></iconify-icon>
                        ${this.formatCreatedDate(resource.createdDate)}
                      </span>
                      <span class="badge badge-info">
                        <iconify-icon icon="icons8:up-round"></iconify-icon>
                        ${resource.seeders}
                      </span>
                    </div>
                    <div class="flex flex-row basis-full justify-end">
                      <button class="btn btn-xs btn-info">Download</button>
                    </div>
                  </div>
                </div>
              `,
            )
          : html`<p>No resources found or loading...</p>`}
      </div>
      ${this.renderPagination()}
    `;
  }
}
