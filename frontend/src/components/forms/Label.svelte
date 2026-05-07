<script lang="ts">
  import type { Snippet } from "svelte";
  import Tooltip from "../interactive/Tooltip.svelte";

  interface Props {
    name: string;
    info?: string;
    ownPositioning?: boolean;
    children?: Snippet;
  }

  let {
    name,
    info = "",
    ownPositioning = true,
    children
  }: Props = $props();
</script>

<style lang="scss">
  @use "../../styles/colors.scss";
  @use "../../styles/dimensions.scss";
  @use "../../styles/text.scss";

  label {
    color: var(--fg-muted);
    font-size: var(--font-size-xs);
    font-weight: var(--font-weight-medium);
    letter-spacing: 0.01em;
    cursor: text;
    display: flex;
    flex-direction: row;
    gap: dimensions.$gapSmaller;
    align-items: center;
    padding-left: 2px;
  }

  .ownPositioning {
    margin-bottom: -(dimensions.$gapMiddle);
  }
</style>

<label for={name} tabindex="-1" class:ownPositioning={ownPositioning}>
  {@render children?.()}
  {#if info}
    <Tooltip tight tiny>
      {info}
    </Tooltip>
  {/if}
</label>