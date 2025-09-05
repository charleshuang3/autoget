import { html, LitElement, unsafeCSS, css, type TemplateResult } from 'lit';
import { customElement, state } from 'lit/decorators.js';
import 'iconify-icon';

import { fetchIndexers, fetchIndexerCategories, type Category } from '../utils/api';
import '../components/navbar.ts';
import globalStyles from '/src/index.css?inline';

@customElement('search-view')
export class SearchView extends LitElement {
  static styles = [
    unsafeCSS(globalStyles),
    css`
      .category-item {
        cursor: pointer;
        transition:
          transform 0.2s,
          box-shadow 0.2s;
      }
      .category-item:hover {
        transform: translateY(-2px);
        box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
      }
      .category-item.active {
        background-color: #4f46e5;
        color: white;
        border-color: #4f46e5;
        transform: scale(1.02);
        box-shadow: 0 6px 12px rgba(0, 0, 0, 0.2);
      }
      .scroll-container {
        overflow-x: auto;
        scrollbar-width: none; /* Firefox */
        -ms-overflow-style: none; /* IE and Edge */
      }
      .scroll-container::-webkit-scrollbar {
        display: none; /* Chrome, Safari, and Opera */
      }
      .break-inside-avoid-column {
        break-inside: avoid-column;
      }
    `,
  ];

  @state()
  private indexers: string[] = [];

  @state()
  private selectedIndexer = '';

  @state()
  private allCategories: Category[] = []; // All categories fetched for the selected indexer

  @state()
  private displayedCategoryLevels: Category[][] = []; // Categories to display for each level

  @state()
  private selectedCategoryPath: Category[] = []; // Path of selected categories

  @state()
  private finalSelectedCategory: Category | null = null; // The leaf category selected

  @state()
  private searchQuery = '';

  async connectedCallback() {
    super.connectedCallback();
    this.indexers = await fetchIndexers();
  }

  private async handleIndexerChange(e: Event) {
    const radio = e.target as HTMLInputElement;
    this.selectedIndexer = radio.value;
    this.selectedCategoryPath = [];
    this.finalSelectedCategory = null;
    await this.fetchAndDisplayCategories();
  }

  private async fetchAndDisplayCategories() {
    if (this.selectedIndexer) {
      this.allCategories = await fetchIndexerCategories(this.selectedIndexer);
      this.displayedCategoryLevels = [this.allCategories]; // Start with the top level
    } else {
      this.allCategories = [];
      this.displayedCategoryLevels = [];
    }
  }

  private handleSearch(e: Event) {
    e.preventDefault();
    console.log({
      query: this.searchQuery,
      indexer: this.selectedIndexer,
      selectedCategory: this.finalSelectedCategory?.id,
      selectedCategoryPath: this.selectedCategoryPath.map((c) => c.id),
    });
  }

  private handleSearchQueryInput(e: Event) {
    this.searchQuery = (e.target as HTMLInputElement).value;
  }

  private handleCategoryClick(category: Category, level: number) {
    // Update the selected path
    this.selectedCategoryPath = this.selectedCategoryPath.slice(0, level);
    this.selectedCategoryPath.push(category);

    // Determine if this is a final selection or if there are subcategories
    if (category.subCategories && category.subCategories.length > 0) {
      this.finalSelectedCategory = null; // Not a final selection yet
      // Add the next level of categories to display
      this.displayedCategoryLevels = this.displayedCategoryLevels.slice(0, level + 1);
      this.displayedCategoryLevels.push(category.subCategories);
    } else {
      this.finalSelectedCategory = category; // This is a leaf category
      // Hide subsequent levels if a leaf node is selected
      this.displayedCategoryLevels = this.displayedCategoryLevels.slice(0, level + 1);
    }
    this.requestUpdate(); // Force re-render to update active states and displayed levels
  }

  private renderCategoryLevel(categories: Category[], level: number): TemplateResult {
    const levelTitle = level === 0 ? 'Main Category' : level === 1 ? 'Sub Category' : 'Tertiary Category';
    const selectedIdInLevel = this.selectedCategoryPath[level]?.id;

    return html`
      <div id="level-${level}" class="flex-shrink-0 w-60 p-2 bg-gray-100 rounded-xl">
        <h3 class="font-semibold text-gray-700 mb-3">${levelTitle}</h3>
        <div id="category-list-${level}" class="flex flex-col space-y-2">
          ${categories.map(
            (category) => html`
              <div
                class="category-item p-2 rounded-lg border border-gray-300 text-left font-medium flex items-center justify-between transition-colors bg-white hover:bg-gray-100 ${selectedIdInLevel ===
                category.id
                  ? 'active'
                  : ''}"
                @click=${() => this.handleCategoryClick(category, level)}
              >
                <span>${category.name}</span>
                ${category.subCategories && category.subCategories.length > 0
                  ? html`<span class="ml-2 text-gray-400 font-bold">›</span>`
                  : ''}
              </div>
            `,
          )}
        </div>
      </div>
    `;
  }

  render() {
    return html`
      <app-navbar activePage="search"></app-navbar>
      <div class="p-4 bg-slate-50 text-gray-800 flex flex-col items-center min-h-screen">
        <div class="bg-white p-8 sm:p-12 rounded-2xl shadow-xl w-full max-w-5xl">
          <form @submit=${this.handleSearch} class="mb-4">
            <div class="join w-full">
              <input
                class="input input-bordered join-item w-full"
                placeholder="Search"
                .value=${this.searchQuery}
                @input=${this.handleSearchQueryInput}
              />
              <button class="btn join-item btn-primary" type="submit">
                <iconify-icon icon="mdi:magnify" width="24" height="24"></iconify-icon>
              </button>
            </div>
          </form>

          <div class="mb-4">
            <p class="text-gray-600 font-medium mb-2">Please choose an indexer:</p>
            <div class="join">
              ${this.indexers.map(
                (indexer) => html`
                  <input
                    class="join-item btn"
                    type="radio"
                    name="indexer"
                    aria-label=${indexer}
                    .value=${indexer}
                    .checked=${this.selectedIndexer === indexer}
                    @change=${this.handleIndexerChange}
                  />
                `,
              )}
            </div>
          </div>

          <p class="text-gray-600 font-medium mb-2">Please choose an category:</p>
          <div class="scroll-container overflow-x-auto flex space-x-4 pb-4">
            ${this.displayedCategoryLevels.map((categories, index) => this.renderCategoryLevel(categories, index))}
          </div>

          ${this.selectedIndexer
            ? html`
                <div class="mt-8 text-center">
                  <p id="selected-path" class="text text-indigo-500 mt-1">
                    ${this.selectedIndexer} > ${this.selectedCategoryPath.map((c) => c.name).join(' > ')}
                  </p>
                </div>
              `
            : ''}
        </div>
      </div>
    `;
  }
}
