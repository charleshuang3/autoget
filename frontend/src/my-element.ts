import { LitElement, html } from "lit";
import { unsafeCSS } from "lit";
import globalStyles from "./index.css?inline";

export class MyElement extends LitElement {
  static styles = [unsafeCSS(globalStyles)];
  render() {
    return html`<input type="checkbox" value="dark" class="toggle theme-controller" /> hello`;
  }
}

window.customElements.define("my-element", MyElement);
