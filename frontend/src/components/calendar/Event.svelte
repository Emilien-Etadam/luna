<script lang="ts">
  import { TextIcon } from "lucide-svelte";

  import { GetEventColor, GetEventHoverColor, GetEventRGB, isDark } from "$lib/common/colors";
  import { passIfEnter } from "$lib/common/inputs";

  import { getContext } from "svelte";
  import type { Writable } from "svelte/store";
  import { NoOp } from "$lib/client/placeholders";
  import ColorCircle from "../misc/ColorCircle.svelte";
  import { getSettings } from "$lib/client/data/settings.svelte";
  import { UserSettingKeys } from "../../types/settings";
  import { getDayIndex } from "$lib/common/date";

  interface Props {
    visible?: boolean;
    event: EventModel | null;
    isFirstDay: boolean;
    date: Date;
    view: "month" | "week" | "day";
    readOnly?: boolean;
  }

  let {
    visible = true,
    event,
    isFirstDay,
    date,
    view,
    readOnly = false,
  }: Props = $props();

  const settings = getSettings();
  let showOnlyCircle = $derived(event && (
    (event.date.allDay && !settings.userSettings[UserSettingKeys.DisplayAllDayEventsFilled]) || 
    (!event.date.allDay && !settings.userSettings[UserSettingKeys.DisplayNonAllDayEventsFilled])
  ));

  let remainingDays = $derived.by(() => {
    if (!date || !event) return 0;
    if (view === "day") return 1;

    const remainingTime = event.date.end.getTime() - date.getTime();
    const remainingDays = Math.ceil(remainingTime / (1000 * 60 * 60 * 24));

    return Math.max(remainingDays, 1);
  })

  let remainingDaysThisWeek = $derived.by(() => {
    const remainingDaysThisWeek = Math.min(remainingDays, 7 - getDayIndex(date));

    return Math.max(remainingDaysThisWeek, 1);
  })

  let eventEndsThisWeek = $derived(remainingDays == remainingDaysThisWeek);

  let currentlyHoveredEvent = $state((getContext("currentlyHoveredEvent") as () => (EventModel | null))());
  let currentlyClickedEvent = $state((getContext("currentlyClickedEvent") as () => (EventModel | null))());

  let showModal: ((event: EventModel) => Promise<EventModel>) = getContext("showEventModal");

  let element: HTMLDivElement | null = $state(null);

  let isEventStart = $derived(event !== null && event.date.start.getTime() >= date.getTime());
  let isFirstDisplay = $derived(isFirstDay || isEventStart);

  let isBackgroundDark: boolean = $derived(event ? isDark(GetEventRGB(event)) : false);
  let participantColors = $derived.by(() => {
    if (!event) return [];
    const own = GetEventColor(event);
    const colors = new Set(event.participant_colors || []);
    if (own) colors.add(own);
    return Array.from(colors);
  });

  function mouseEnter() {
    if (event == null) return;

    currentlyHoveredEvent = event;
  }
  function mouseLeave() {
    if (event == null) return;

    if (currentlyHoveredEvent == event)
      currentlyHoveredEvent = null;
    if (currentlyClickedEvent == event)
      currentlyClickedEvent = null;
  }
  function mouseDown() {
    if (event == null) return;

    currentlyClickedEvent = event;
  }
  function mouseUp() {
    if (event == null) return;

    if (currentlyClickedEvent == event) {
      currentlyClickedEvent = null;
      showModal(event).then(newEvent => event = newEvent).catch(NoOp);
      element?.blur();
    }
  }
  function keyPress(e: KeyboardEvent) {
    passIfEnter(e, () => {
      if (event) showModal(event).then(newEvent => event = newEvent).catch(NoOp);
      element?.blur();
    });
  }
</script>

<style lang="scss">
  @use "../../styles/animations.scss";
  @use "../../styles/colors.scss";
  @use "../../styles/dimensions.scss";
  @use "../../styles/text.scss";

  div {
    box-sizing: border-box;
    min-height: 18px;
    padding: 2px 4px;
    padding-left: calc(var(--gapBetweenDays) + 4px);
    font-size: var(--font-size-ui);
    margin: 0;

    display: flex;
    gap: 6px;
    flex-direction: row;
    flex-wrap: nowrap;
    align-items: center;

    user-select: none;
    cursor: pointer;

    white-space: nowrap;
    overflow: visible;

    flex-shrink: 0;

    border-left: 2px solid transparent;
    transition: background-color linear animations.$animationSpeedFast, border-color linear animations.$animationSpeedFast;
  }

  div:focus {
    outline: none;
  }

  div::after {
    content: ".";
    visibility: hidden;
  }
  div.placeholder {
    visibility: hidden;
  }
  div.start {
    border-top-left-radius: var(--radius-2);
    border-bottom-left-radius: var(--radius-2);
    margin-left: var(--gapBetweenDays);
    padding-left: 4px;
  }
  div.end {
    border-top-right-radius: var(--radius-2);
    border-bottom-right-radius: var(--radius-2);
    margin-right: var(--gapBetweenDays);
  }

  div.hidden {
    display: none;
  }

  div.foregroundBright {
    color: colors.$foregroundBright;
  }
  div.foregroundDark {
    color: colors.$foregroundDark;
  }

  span.name {
    text-overflow: ellipsis;
    overflow: hidden;
    min-width: 0;
    flex-shrink: 1;
    font-size: var(--font-size-ui);
    font-weight: 400;
  }
  span.time {
    flex-shrink: 0;
    text-align: center;
    font-weight: 400;
    font-family: text.$fontFamilyTime;
    font-size: var(--font-size-event-time);
    color: var(--fg-muted);
  }
  span.icons {
    flex-shrink: 0;
    display: flex;
    align-items: center;
  }
  span.participants {
    flex-shrink: 0;
    display: flex;
    align-items: center;
    gap: 4px;
  }

  div.onlyCircle {
    background-color: transparent !important;
    color: var(--fg-muted) !important;
  }

  div.hover {
    border-left-color: transparent;
    background-color: var(--bg-hover) !important;
    color: var(--fg-primary) !important;
  }

  div.active {
    border-left-color: var(--border-focus);
    background-color: var(--bg-selection-active) !important;
    color: var(--fg-strong) !important;
  }
</style>

<!-- TODO: the following reduced the amount of divs we need to render but was prone to some edge-case bugs (no.116) -->
<!--{#if event && (isFirstDisplay || getDayIndex(date) == 0 || showOnlyCircle)}-->
{#if event}
  <div
    bind:this={element}
    class:start={isEventStart}
    class:end={eventEndsThisWeek}
    class:hover={currentlyHoveredEvent == event}
    class:active={currentlyClickedEvent == event}
    class:hidden={!visible}
    class:foregroundBright={isBackgroundDark}
    class:foregroundDark={!isBackgroundDark}
    class:onlyCircle={showOnlyCircle}
    onmouseenter={mouseEnter}
    onmouseleave={mouseLeave}
    onmousedown={mouseDown}
    onmouseup={mouseUp}
    onfocusin={mouseEnter}
    onfocusout={mouseLeave}
    onkeypress={keyPress}
    role="button"
    tabindex={isFirstDisplay ? 0 : -1}
    style="
      background-color:{currentlyHoveredEvent == event ? GetEventHoverColor(event) : GetEventColor(event)};
      width: calc({(showOnlyCircle ? 1 : remainingDaysThisWeek) * 100}% - {((isEventStart ? 1 : 0) + (eventEndsThisWeek ? 1 : 0)) * (showOnlyCircle ? 0 : 1)} * var(--gapBetweenDays));
      z-index: {16 - getDayIndex(date)};
    "
  >
    {#if showOnlyCircle}
      <ColorCircle
        color={GetEventColor(event)}
        size="small"
      />
    {/if}
    {#if !event.date.allDay && event.date.start >= date}
      <span class="time">
        {event.date.start.toLocaleTimeString([], {hour: '2-digit', minute: '2-digit'})}
      </span>
    {/if}
    <span class="name">
      {event.name}
    </span>
    {#if !showOnlyCircle && participantColors.length > 1}
      <span class="participants">
        {#each participantColors as color (color)}
          <ColorCircle color={color} size="small" />
        {/each}
      </span>
    {/if}
    {#if (event.desc && event.desc != "")}
      <span class="icons">
        <TextIcon size={12}/>
      </span>
    {/if}
  </div>
{:else}
  <div class="placeholder" class:hidden={!visible}></div>
{/if}
