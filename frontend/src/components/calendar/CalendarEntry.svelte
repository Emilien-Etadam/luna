<script lang="ts">
  import ColorCircle from "../misc/ColorCircle.svelte";
  import Spinner from "../decoration/Spinner.svelte";
  import Tooltip from "../interactive/Tooltip.svelte";
  import VisibilityToggle from "../interactive/VisibilityToggle.svelte";

  import { GetCalendarColor } from "$lib/common/colors";
  import { NoOp } from "$lib/client/placeholders";
  import { focusIndicator } from "$lib/client/decoration";
  import { getMetadata } from "$lib/client/data/metadata.svelte";
  import { draggable } from "$lib/client/reordering.svelte";

  import { getContext } from "svelte";
  import { queueNotification } from "../../lib/client/notifications";
  import { ColorKeys } from "../../types/colors";
  import { getRepository } from "../../lib/client/data/repository.svelte";

  interface Props {
    calendar: CalendarModel;
    readOnly?: boolean;
  }

  let { calendar = $bindable(), readOnly = false }: Props = $props();

  const metadata = getMetadata();
  const repository = getRepository();

  let hasErrored = $derived(calendar && metadata.faultyCalendars.has(calendar.id));
  let isLoading = $derived(calendar && metadata.loadingCalendars.get(calendar.id));
  let calendarVisible = $derived(calendar && metadata.hiddenCalendars.has(calendar.id));

  let showModal: ((calendar: CalendarModel) => Promise<CalendarModel>) = getContext("showCalendarModal");
  function showModalInternal() {
    showModal(calendar).then(newCalendar => calendar = newCalendar).catch(NoOp);
  }

  $effect(() => {
    const shouldBeVisible = !metadata.hiddenCalendars.has(calendar.id);
    if (shouldBeVisible == calendarVisible) return;
    calendarVisible = shouldBeVisible;
  });
  function setVisible(visible: boolean) {
    getMetadata().setCalendarVisibility(calendar.id, visible);
  }

  async function reorderCalendar(newIndex: number) {
    if (readOnly) return;
    await repository.changeCalendarDisplayOrder(calendar, newIndex).catch((err) => {
      queueNotification(ColorKeys.Danger, err);
    });
  }
</script>

<style lang="scss">
  @use "../../styles/colors.scss";
  @use "../../styles/dimensions.scss";

  div.calendarEntry {
    display: flex;
    flex-direction: row;
    gap: 8px;
    width: 100%;
    align-items: center;
    justify-content: space-between;
    user-select: none;
    cursor: grab;
    min-height: 24px;
    padding: 2px 8px 2px 14px;
    border-radius: var(--radius-1);
    transition: background-color var(--transition-fast);
  }

  div.calendarEntry:hover {
    background-color: var(--bg-hover);
  }

  div.calendarEntry:focus-within {
    background-color: var(--bg-hover);
  }

  span {
    display: flex;
    flex-direction: row;
    align-items: center;
  }

  span.name {
    gap: dimensions.$gapSmall;
    min-width: 0;
  }

  span.buttons {
    gap: dimensions.$gapTiny;
  }

  button {
    all: unset;
    cursor: pointer;
    display: inline;
    width: max-content;
    position: relative;
    text-wrap: nowrap;
    text-overflow: ellipsis;
    min-width: 0;
    overflow: hidden;
  }

  span.calendarName {
    font-size: var(--font-size-ui);
    font-weight: 400;
    color: var(--fg-primary);
    text-wrap: nowrap;
    text-overflow: ellipsis;
    min-width: 0;
    overflow: hidden;
  }
</style>

{#if readOnly}
  <div class="calendarEntry">
    <span class="name">
      <ColorCircle
        color={GetCalendarColor(calendar)}
        size="small"
      />
      <span class="calendarName">{calendar.name}</span>
    </span>
    <span class="buttons">
      {#if isLoading}
        <Spinner/>
      {/if}
      <VisibilityToggle bind:visible={calendarVisible} onClick={setVisible}/>
      {#if hasErrored}
        <Tooltip error={true}>An error occurred trying to retrieve events from this calendar.</Tooltip>
      {/if}
    </span>
  </div>
{:else}
  <div class="calendarEntry" use:draggable={{ ownClass: "calendarEntry", childClasses: [], callback: reorderCalendar}}>
    <span class="name">
      <ColorCircle
        color={GetCalendarColor(calendar)}
        size="small"
      />
      <button onclick={showModalInternal} use:focusIndicator={{ type: "underline" }}>
        {calendar.name}
      </button>
    </span>
    <span class="buttons">
      {#if isLoading}
        <Spinner/>
      {/if}
      <VisibilityToggle bind:visible={calendarVisible} onClick={setVisible}/>
      {#if hasErrored}
        <Tooltip error={true}>An error occurred trying to retrieve events from this calendar.</Tooltip>
      {/if}
    </span>
  </div>
{/if}