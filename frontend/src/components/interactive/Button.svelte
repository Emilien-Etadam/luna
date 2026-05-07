<script lang="ts">
  import type { Snippet } from "svelte";
  import { ColorKeys } from "../../types/colors";
  import { addRipple, focusIndicator } from "../../lib/client/decoration";

  interface Props {
    onClick?: () => void;
    color?: ColorKeys;
    type?: "button" | "submit";
    compact?: boolean;
    enabled?: boolean;
    href?: string;
    children?: Snippet;
  }

  let {
    onClick = () => {},
    color = ColorKeys.Neutral,
    type = "button",
    compact = false,
    enabled = true,
    href = "",
    children
  }: Props = $props();
</script>

<style lang="scss">
  @use "sass:map";

  @use "../../styles/animations.scss";
  @use "../../styles/colors.scss";
  @use "../../styles/dimensions.scss";
  @use "../../styles/text.scss";

  button, a {
    /* reset */
    background: none;
    color: inherit;
    border: 1px solid transparent;
    padding: 0;
    font: inherit;
    cursor: pointer;
    outline: inherit;
    text-decoration: none;

    display: inline-flex;
    align-items: center;
    justify-content: center;

    cursor: pointer;
    padding: 0 dimensions.$gapSmall;
    height: var(--control-h-md);
    line-height: 1;
    border-radius: var(--borderRadius);

    min-width: dimensions.$buttonMinWidth;
    text-align: center;
    font-weight: var(--font-weight-medium);
    letter-spacing: 0.005em;
    
    position: relative;
    overflow: hidden;

    transition: background-color var(--transition-fast), color var(--transition-fast),
                border-color var(--transition-fast), filter var(--transition-fast);
  }

  /* Hack pour aligner le spinner sur la baseline du bouton */
  button > :global(span.spinner) {
    &::before, &::after {
      content: "a";
      visibility: hidden;
    }
  }

  button:not(.neutral) {
    --barFocusIndicatorColor: #{colors.$barFocusIndicatorColorAlt};
  }

  button.compact, a.compact {
    min-width: dimensions.$buttonMinWidthCompact;
    height: var(--control-h-sm);
    padding: 0 dimensions.$gapSmaller;
    font-size: var(--font-size-sm);
  }

  button:focus-visible, a:focus-visible {
    box-shadow: var(--focus-ring);
  }

  /* Hover/active subtils, sans écraser la couleur sémantique */
  button:not(.disabled):hover, a:not(.disabled):hover {
    filter: brightness(1.08);
  }
  button:not(.disabled):active, a:not(.disabled):active {
    filter: brightness(0.94);
  }

  .disabled {
    cursor: not-allowed;
    opacity: 0.6;
  }

  /* Variantes de couleur (utilisent les mêmes alias colors.$specialColors) */
  @each $key, $val in colors.$specialColors {
    button.#{$key}, a.#{$key} {
      background-color: map.get($val, "background");
      color: map.get($val, "foreground");
    }
    button.#{$key}.disabled, a.#{$key}.disabled {
      color: color-mix(in srgb, map.get($val, "foreground") 50%, transparent);
    }
  }

  /* Neutral: bouton "fantôme" avec bordure subtile pour rester visible
     même posé sur la sidebar (bg identique). */
  button.neutral, a.neutral {
    background-color: transparent;
    color: var(--fg-primary);
    border-color: var(--border-default);
  }
  button.neutral:hover, a.neutral:hover {
    background-color: var(--bg-hover);
    color: var(--fg-strong);
    filter: none;
  }
  button.neutral:active, a.neutral:active {
    background-color: var(--bg-active);
    filter: none;
  }
</style>

{#if href !== ""}
  <a
    class:success={color == ColorKeys.Success}
    class:warning={color == ColorKeys.Warning}
    class:danger={color == ColorKeys.Danger}
    class:accent={color == ColorKeys.Accent}
    class:neutral={color == ColorKeys.Neutral}
    class:inherit={color == ColorKeys.Inherit}
    class:compact={compact}
    onmouseleave={(e) => {(e.target as HTMLButtonElement).blur()}}
    class:disabled={!enabled}
    href={enabled ? href : "#"}
    onmousedown={addRipple}
    use:focusIndicator
  >
    {@render children?.()}
  </a>
{:else}
  <button
    class:success={color == ColorKeys.Success}
    class:warning={color == ColorKeys.Warning}
    class:danger={color == ColorKeys.Danger}
    class:accent={color == ColorKeys.Accent}
    class:neutral={color == ColorKeys.Neutral}
    class:inherit={color == ColorKeys.Inherit}
    class:compact={compact}
    onclick={onClick}
    onmouseleave={(e) => {(e.target as HTMLButtonElement).blur()}}
    type={type}
    disabled={!enabled}
    class:disabled={!enabled}
    onmousedown={addRipple}
    use:focusIndicator
  >
    {@render children?.()}
  </button>
{/if}