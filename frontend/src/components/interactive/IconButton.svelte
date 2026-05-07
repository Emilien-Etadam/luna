<script lang="ts">
  import type { Snippet } from "svelte";

  import { NoOp } from "$lib/client/placeholders";
  import Tooltip from "./Tooltip.svelte";

  interface Props {
    up?: () => void;
    down?: () => void;
    click?: () => void;
    visible?: boolean;
    info?: string;
    style?: string;
    tabindex?: number;
    href?: string;
    children?: Snippet;
  }

  let {
    up = NoOp,
    down = NoOp,
    click = NoOp,
    visible = true,
    info = "",
    style = "",
    tabindex = 0,
    href = "",
    children
  }: Props = $props();

  // svelte-ignore non_reactive_update
  // isLink is set once and never changed
  let button = $state<HTMLElement | null>(null);

  function clickInternal(e: MouseEvent) {
    e.stopPropagation();
    click();
  }

  function leaveInternal() {
    if (!button) return;
    button.blur();
    up();
  }
  function upInternal() {
    if (!button) return;
    button.blur();
    up();
  }
</script>

<style lang="scss">
  @use "../../styles/animations.scss";
  @use "../../styles/colors.scss";
  @use "../../styles/dimensions.scss";

  button, a {
    all: unset;
    border-radius: var(--borderRadius);
    display: flex;
    align-items: center;
    justify-content: center;
    height: var(--control-h-md);
    min-width: var(--control-h-md);
    padding: 0 6px;
    cursor: pointer;
    position: relative;
    color: inherit;
    transition: background-color var(--transition-fast), color var(--transition-fast);
  }

  button.hidden, a.hidden {
    visibility: hidden;
  }

  /* Surface de hover : remplie via ::before pour ne pas dépendre d'une div interne */
  button::before, a::before {
    content: "";
    position: absolute;
    inset: 0;
    background-color: var(--bg-hover);
    border-radius: inherit;
    opacity: 0;
    transition: opacity var(--transition-fast);
    pointer-events: none;
    z-index: -1;
  }

  button:hover::before, a:hover::before,
  button:focus-visible::before, a:focus-visible::before {
    opacity: 1;
  }

  button:active::before, a:active::before {
    opacity: 1;
    background-color: var(--bg-active);
  }

  button:focus-visible, a:focus-visible {
    box-shadow: var(--focus-ring);
  }

  button:hover, button:focus-visible, a:hover, a:focus-visible {
    color: var(--fg-primary);
  }

  /* On garde div.circle pour compat (IconButton existant), mais on
     l'invisibilise ; les états sont gérés via ::before */
  div.circle { display: none; }
</style>

{#if info == ""}
  {@render buttonSnippet()}
{:else}
  <Tooltip
    icon={buttonSnippet} 
    inheritColor={true}
    tight={true}
    pointerCursor={true}
  >
    {info}
  </Tooltip>
{/if}

{#snippet buttonSnippet()}
  {#if href !== ""}
    <a
      bind:this={button}
      class:hidden={!visible}
      href={href}
      style={style}
      tabindex="{tabindex}"
    >
      <div class="circle"></div>
      {@render children?.()}
    </a>
  {:else}
    <button
      bind:this={button}
      onclick={clickInternal}
      onmousedown={down}
      onmouseleave={leaveInternal}
      onmouseup={upInternal}
      class:hidden={!visible}
      type="button"
      style={style}
      tabindex="{tabindex}"
    >
      <div class="circle"></div>
      {@render children?.()}
    </button>
  {/if}
{/snippet}