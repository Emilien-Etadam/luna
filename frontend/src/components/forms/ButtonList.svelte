<script lang="ts" generics="T">
  import { addRipple, focusIndicator } from "../../lib/client/decoration";
  import { EmptyOption } from "../../lib/client/placeholders";
  import { ColorKeys } from "../../types/colors";
  import type { Option } from "../../types/options";

  // This component is used for category lists (currently only in the settings modal)

  interface Props {
    value: T;
    options: Option<T>[][];
  }

  let {
    value = $bindable(),
    options
  }: Props = $props();

  let selected: Option<T> = $derived(options.flat().filter(option => option.value === value)[0] || options[0] || EmptyOption);
</script>

<style lang="scss">
  @use "sass:map";

  @use "../../styles/colors.scss";
  @use "../../styles/dimensions.scss";

  aside {
    display: flex;
    flex-direction: column;
    gap: 1px;
    overflow: auto;
    padding-right: dimensions.$gapSmall;
    margin-right: -(dimensions.$gapSmall);
    min-width: 11em;
  }

  .option {
    border: 0;
    outline: 0;
    font-family: inherit;
    font-size: inherit;
    width: 100%;
    background-color: transparent;
    color: var(--fg-muted);
    padding: 7px 10px;
    border-radius: var(--radius-1);
    display: flex;
    flex-direction: row;
    justify-content: flex-start;
    align-items: center;
    gap: dimensions.$gapSmall;
    position: relative;
    overflow: hidden;
    cursor: pointer;
    flex-shrink: 0;
    transition: background-color var(--transition-fast), color var(--transition-fast);
  }

  .option:hover {
    background-color: var(--bg-hover);
    color: var(--fg-primary);
  }

  .option:focus-visible {
    box-shadow: var(--focus-ring);
  }

  .first {
    margin-top: 0;
  }

  :not(.option:last-child).last {
    margin-bottom: dimensions.$gapSmall;
  }

  /* Sélection : fond subtle + accent à gauche pour rester sobre */
  .option.selected {
    background-color: var(--bg-active);
    color: var(--fg-strong);
  }
  .option.selected.danger {
    color: var(--accent-red);
  }
  .option.selected.warning {
    color: var(--accent-yellow);
  }

  .option :global(*) {
    pointer-events: none;
  }

  p {
    margin: 0;
    padding: 0;
    flex-grow: 1;
    text-align: left;
    font-size: var(--font-size-ui);
    font-weight: var(--font-weight-medium);
  }
</style>

<aside>
  {#each options as block}
    {#each block as option, i}
      {@const Icon = option.icon}
      <button
        class="option"
        class:first={i === 0}
        class:last={i === block.length - 1}
        class:selected={option.value === value}
        class:success={option.color === ColorKeys.Success}
        class:warning={option.color === ColorKeys.Warning}
        class:danger={option.color === ColorKeys.Danger}
        class:accent={!option.color || option.color === ColorKeys.Accent}
        class:neutral={option.color === ColorKeys.Neutral}
        class:inherit={option.color === ColorKeys.Inherit}
        onclick={() => value = option.value}
        onmousedown={addRipple}
        use:focusIndicator
      >
        <Icon size={20}/>
        <p>
          {option.name}
        </p>
      </button>
    {/each}
  {/each}
</aside>