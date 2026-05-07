<script lang="ts">
  import { focusIndicator } from "$lib/client/decoration";

  interface Props {
    value: boolean;
    name: string;
    enabled?: boolean;
    onChange?: (value: boolean) => any;
    toggle: (e: MouseEvent | KeyboardEvent) => void;
  }

  let {
    value = $bindable(),
    name,
    enabled = true,
    onChange = () => {},
    toggle = $bindable(),
  }: Props = $props();

  toggle = (e: MouseEvent | KeyboardEvent) => {
    value = !value;
    onChange(value);
    e.stopPropagation();
  }
</script>

<style lang="scss">
  @use "../../styles/animations.scss";
  @use "../../styles/colors.scss";
  @use "../../styles/dimensions.scss";
  @use "../../styles/text.scss";

  button {
    all: unset;

    width: calc(2 * text.$lineHeightParagraph);
    height: text.$lineHeightParagraph;

    position: relative;

    display: flex;
    justify-content: center;
    align-items: center;

    cursor: pointer;

    border-radius: var(--radius-pill);
    background-color: var(--surface-overlay);
    border: 1px solid var(--border-default);
    overflow: hidden;
    transition: background-color var(--transition-fast), border-color var(--transition-fast),
                box-shadow var(--transition-fast);
  }

  button:hover {
    border-color: var(--border-strong);
  }

  button:focus-visible {
    box-shadow: var(--focus-ring);
  }

  button.check {
    background-color: var(--accent-blue);
    border-color: var(--accent-blue);
  }

  button.check {
    --barFocusIndicatorColor: #{colors.$barFocusIndicatorColorAlt} !important;
  }

  button::after {
    content: "";
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background-color: transparent;
    transition: opacity var(--transition-fast);
    opacity: 0;
  }

  .handle {
    height: calc(100% - 6px);
    aspect-ratio: 1/1;
    border-radius: 50%;
    background-color: var(--fg-muted);

    left: 3px;

    position: absolute;

    transition: transform var(--transition) var(--easing-standard),
                background-color var(--transition-fast);

    z-index: 2;
  }

  button.check .handle {
    background-color: #ffffff;
  }

  .handle.check {
    transform: translateX(calc(100% + 3px));
  }

  button.disabled {
    cursor: unset;
  }

  button :global(*) {
    pointer-events: none;
  }

  input {
    all: unset;
    position: absolute;
    left: 0;
    top: 0;
  }
</style>

<!-- Components that use this toggle all implement for={name} -->
<!-- svelte-ignore a11y_consider_explicit_label -->
<button
  type="button"
  class:disabled={!enabled}
  class:check={value}
  onclick={toggle}
  use:focusIndicator
>
  <div
    class="handle"
    class:check={value}
  >
  </div>
  <input type="hidden" name={name} value={value}>
</button>