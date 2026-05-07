<script lang="ts">
  import { PlusIcon } from "lucide-svelte";
  import { getContext } from "svelte";

  import Event from "./Event.svelte";
  import IconButton from "../interactive/IconButton.svelte";

  import { queueNotification } from "$lib/client/notifications";
  import { NoOp } from "$lib/client/placeholders";
  import { ColorKeys } from "../../types/colors";

  interface Props {
    date: Date;
    isCurrentMonth: boolean;
    isFirstDay: boolean;
    isToday: boolean;
    events: (EventModel | null)[];
    maxEvents?: number;
    containerHeight: number;
    view: "month" | "week" | "day";
    showMore?: (date: Date, events: (EventModel | null)[]) => any;
    readOnly?: boolean;
  }

  let {
    date,
    isCurrentMonth,
    isFirstDay,
    isToday,
    events,
    maxEvents = 1,
    containerHeight = $bindable(),
    view,
    showMore = NoOp,
    readOnly = false,
  }: Props = $props();

  let showCreateEventModal: ((date: Date) => Promise<EventModel>) = getContext("showNewEventModal");
  let createEventButtonClick = () => {
    showCreateEventModal(date).catch((err) => {
      if (err) queueNotification(ColorKeys.Danger, `Could not create event: ${err.message}`);
    });
  };

  let actualMaxEvents: number = $derived(maxEvents <= events.length - 1 ? maxEvents - 1 : maxEvents);
</script>

<style lang="scss">
  @use "../../styles/animations.scss";
  @use "../../styles/colors.scss";
  @use "../../styles/dimensions.scss";
  @use "../../styles/text.scss";

  div.day {
    min-width: 0;
    overflow: visible;
    height: 100%;
    position: relative;
    font-size: var(--font-size-day-number);
    --gapBetweenDays: 0px;
    border-right: 1px solid var(--border-subtle);
    border-bottom: 1px solid var(--border-subtle);
    transition: background-color var(--transition-fast);
  }

  /* Suppression de la double bordure au bord droit/bas de la grille
     (le wrapper main fournit déjà une bordure) */
  :global(div.days > div.day:nth-child(7n)) {
    border-right: 0;
  }

  div.background {
    display: flex;
    flex-direction: column;
    gap: 4px;
    margin: 0;
    padding: 6px 8px;
    background-color: var(--bg-editor);
    border: 0;
    height: 100%;
    transition: background-color var(--transition-fast);
  }

  div.background.weekend {
    background-color: var(--bg-weekend);
  }

  div.background.otherMonth {
    background-color: var(--bg-other-month);
  }

  div.background.today {
    background-color: var(--bg-today);
  }

  /* Hover doux sur la cellule entière (uniquement vue mois) */
  div.day:hover div.background:not(.today) {
    background-color: color-mix(in srgb, var(--tint-hover) 30%, var(--bg-editor));
  }
  div.day:hover div.background.weekend:not(.today) {
    background-color: color-mix(in srgb, var(--tint-hover) 30%, var(--bg-weekend));
  }
  div.day:hover div.background.otherMonth:not(.today) {
    background-color: color-mix(in srgb, var(--tint-hover) 30%, var(--bg-other-month));
  }

  span.top {
    display: grid;
    grid-template-columns: 1fr 1fr 1fr;
    grid-template-areas: "none date add";
    align-items: center;
    min-height: 18px;
    font-size: var(--font-size-day-number);
  }

  span.date {
    position: relative;
    z-index: 0;
    color: var(--fg-primary);
    font-size: var(--font-size-day-number);
    font-variant-numeric: tabular-nums;
    font-weight: 400;
    text-align: center;
    width: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    grid-area: date;
    user-select: none;
    line-height: 1;
  }
  span.date.weekend {
    color: var(--fg-muted);
  }
  span.date.otherMonth {
    color: var(--fg-other-month);
  }

  span.date.today {
    color: var(--fg-today);
    font-weight: var(--font-weight-semibold);
  }
  span.date.today::before {
    content: "";
    position: absolute;
    z-index: -1;
    left: 50%;
    top: 50%;
    transform: translate(-50%, -50%);
    width: 20px;
    height: 20px;
    border-radius: var(--radius-pill);
    background-color: var(--accent-blue);
    box-shadow: 0 0 0 1px color-mix(in srgb, var(--accent-blue) 35%, transparent);
  }

  span.add {
    grid-area: add;
    display: flex;
    align-items: center;
    justify-content: right;
    opacity: 0;
    transition: opacity var(--transition-fast);
  }
  div.day:hover span.add,
  div.day:focus-within span.add {
    opacity: 1;
  }

  button.more {
    all: unset;
    text-align: left;
    color: var(--fg-muted);
    background-color: transparent;
    cursor: pointer;
    z-index: 20;
    margin: 2px var(--gapBetweenDays) 0;
    padding: 2px 6px;
    border-radius: var(--radius-1);
    position: relative;
    font-size: var(--font-size-xs);
    font-weight: var(--font-weight-medium);
    line-height: 1.2;
    transition: background-color var(--transition-fast), color var(--transition-fast);
  }

  button.more:hover,
  button.more:focus-visible {
    background-color: var(--bg-hover);
    color: var(--fg-primary);
  }

  button.more.otherMonth {
    color: var(--fg-other-month);
  }

  div.events {
    position: absolute;
    display: flex;
    flex-direction: column;
    gap: 3px;

    --topMargin: calc(var(--font-size-day-number) + 16px);
    top: var(--topMargin);
    height: calc(100% - var(--topMargin) - var(--gapBetweenDays));
    width: 100%;
  }
</style>

<div class="day">
  <div
    class="background"
    class:otherMonth={!isCurrentMonth}
    class:weekend={(date.getDay() === 0 || date.getDay() === 6) && isCurrentMonth}
    class:today={isToday}
  >
    <span class="top">
      <span
        class="date"
        class:weekend={date.getDay() === 0 || date.getDay() === 6}
        class:today={isToday}
        class:otherMonth={!isCurrentMonth}
      >
        {date.getDate()}
      </span>
      {#if !readOnly}
        <span class="add">
          <IconButton click={createEventButtonClick} tabindex={-1}>
            <PlusIcon size={13}/>
          </IconButton>
        </span>
      {/if}
    </span>
  </div>
  {#if isFirstDay}
    <div class="events" bind:offsetHeight={containerHeight}>
      {@render eventEntries()}
    </div>
  {:else}
    <div class="events">
      {@render eventEntries()}
    </div>
  {/if}
</div>

{#snippet eventEntries()}
  <!-- TODO: forcing EventEntry to be unique for each event and i like that
  fixes a few issues but might be less performant. figure out the right
  compromise -->
  <!-- {#each events as event, i ((event?.id || 0) + i.toString())} -->

  {#each events as event, i ((event?.id || i).toString() + date.getTime())}
    <!-- TODO: make parameters match css, look into cubic easing, invert fly direction when going back in range -->
    <!--<div
      class="eventAnimation"
      animate:flip={{duration: 300, delay: 300}}
      in:fly={{duration: 300, x: 200}}
      out:fly={{duration: 300, x: -200}}
      style="z-index: {16 - getDayIndex(date)}"
      class:hidden={i >= actualMaxEvents}
    >-->
      <Event
        event={event}
        isFirstDay={isFirstDay}
        date={date}
        visible={i < actualMaxEvents}
        view={view}
        {readOnly}
      />
    <!--</div>-->
  {/each}
  {#if events.length > maxEvents && actualMaxEvents >= 0}
    <button class="more" class:otherMonth={!isCurrentMonth} onclick={() => showMore(date, events)}>
      {#if actualMaxEvents == 0}
       {events.length} events
      {:else}
        and {events.length - actualMaxEvents} more
      {/if}
    </button>
  {/if}
{/snippet}