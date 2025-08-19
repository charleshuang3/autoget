import { LitElement, html, unsafeCSS } from 'lit';
import { customElement, state } from 'lit/decorators.js';
import { until } from 'lit/directives/until.js';
import globalStyles from '/src/index.css?inline';

@customElement('app-navbar')
export class AppNavbar extends LitElement {
  static styles = [unsafeCSS(globalStyles)];

  @state()
  private _indexers: string[] = [];

  connectedCallback() {
    super.connectedCallback();
    this._fetchIndexers();
  }

  private async _fetchIndexers() {
    try {
      const response = await fetch('/api/v1/indexers');
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      this._indexers = await response.json();
    } catch (error) {
      console.error('Failed to fetch indexers:', error);
      this._indexers = []; // Set to empty array on error
    }
  }

  render() {
    return html`
      <div class="navbar bg-base-200">
        <div class="navbar-start">
          <a href="/" class="btn-ghost">
            <img src="/icon.svg" alt="Icon" class="w-8 h-8" />
          </a>
          <div role="tablist" class="tabs tabs-border">
            ${until(
              this._indexers.map((indexer) => {
                return html`<a href="/indexers/${indexer}" class="tab" role="tab">${indexer}</a>`;
              }),
            )}
          </div>
        </div>
        <div class="navbar-end">
          <a href="/search" class="btn btn-ghost">Search</a>
        </div>
      </div>
    `;
  }
}
