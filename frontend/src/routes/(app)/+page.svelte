<script lang="ts">
  import { Github } from "svelte-simples"
  import { CalendarDays, CalendarRange, Calendar as CalendarSingle, Copyleft, List, PlusIcon, RefreshCw, Settings, WifiOff } from "lucide-svelte";
  import { onMount, setContext, untrack } from "svelte";

  import Calendar from "../../components/calendar/Calendar.svelte";
  import CalendarEntry from "../../components/calendar/CalendarEntry.svelte";
  import CalendarModal from "../../components/modals/CalendarModal.svelte";
  import EventModal from "../../components/modals/EventModal.svelte";
  import Horizontal from "../../components/layout/Horizontal.svelte";
  import IconButton from "../../components/interactive/IconButton.svelte";
  import MonthSelection from "../../components/interactive/MonthSelection.svelte";
  import SourceEntry from "../../components/calendar/SourceEntry.svelte";
  import SourceModal from "../../components/modals/SourceModal.svelte";

  import { afterNavigate, beforeNavigate } from "$app/navigation";
  import { browser } from "$app/environment";

  import SmallCalendar from "../../components/interactive/SmallCalendar.svelte";
  import { NoOp } from "$lib/client/placeholders";
  import { getMetadata } from "$lib/client/data/metadata.svelte";
  import { getRepository } from "$lib/client/data/repository.svelte";
  import { queueNotification } from "$lib/client/notifications";
  import { getConnectivity, Reachability } from "$lib/client/data/connectivity.svelte";
  import Button from "../../components/interactive/Button.svelte";
  import DayViewModal from "../../components/modals/DayViewModal.svelte";
  import { getDayIndex, isInRange } from "../../lib/common/date";
  import { compareEventsByStartDate } from "../../lib/common/comparators";
  import { GetEventColor } from "$lib/common/colors";
  import SourceWizardModal from "../../components/modals/SourceWizardModal.svelte";
  import SettingsModal from "../../components/modals/SettingsModal.svelte";
  import { getSettings } from "$lib/client/data/settings.svelte";
  import { UserSettingKeys } from "../../types/settings";
  import ThemeToggle from "../../components/interactive/ThemeToggle.svelte";
  import { ColorKeys } from "../../types/colors";
  import { page } from "$app/state";
  import CreditsModal from "../../components/modals/CreditsModal.svelte";

  /* Singletons */
  const settings = getSettings();
  const metadata = getMetadata();
  const repository = getRepository();
  const connectivity = getConnectivity();
  const publicReadonly = $derived(!settings.userData.id);

  /* Constants */
  let autoRefreshInterval = 1000 * 60; // 1 minute

  /* View logic */
  let view: "month" | "week" | "day" | "agenda" = $derived.by(() => {
    const stored = page.url.searchParams.get("view");
    if (!stored || !["month", "week", "day", "agenda"].includes(stored)) return "month"
    return stored as "month" | "week" | "day" | "agenda";
  });

  let today = $state(new Date());
  let date = $derived.by(() => {
    const stored = page.url.searchParams.get("date");
    if (!stored) return new Date(today);
    const parsed = new Date(stored);
    return parsed;
  });

  $effect(() => {
    const url = new URL(window.location.toString());

    url.searchParams.set("view", view);
    url.searchParams.set("date", date.toISOString().split("T")[0]);

    history.replaceState(history.state, '', url);
  })

  function getVisibleRange(date: Date, view: "month" | "week" | "day" | "agenda"): { start: Date, end: Date } {
    const rangeStart = new Date(date);
    const rangeEnd = new Date(date);
    rangeStart.setHours(0, 0, 0, 0);
    rangeEnd.setHours(23, 59, 59, 999);
    switch (view) {
      case "month":
        rangeStart.setDate(1);
        rangeStart.setDate(rangeStart.getDate() - getDayIndex(rangeStart));
        rangeEnd.setMonth(rangeEnd.getMonth() + 1);
        rangeEnd.setDate(0);
        rangeEnd.setDate(rangeEnd.getDate() + 6 - getDayIndex(rangeEnd));
        break;
      case "week":
        rangeStart.setDate(date.getDate() - getDayIndex(rangeStart));
        rangeEnd.setDate(rangeStart.getDate() + 7);
        break;
      case "day":
      case "agenda":
        rangeStart.setDate(1);
        rangeEnd.setMonth(rangeEnd.getMonth() + 1);
        rangeEnd.setDate(0);
        break;
      default:
    }
    return { start: rangeStart, end: rangeEnd };
  }
  let todayInRange = $derived.by(() => {
    const range = getVisibleRange(date, view);
    return isInRange(today, range.start, range.end)
  });

  function seeToday() {
    today = new Date();
    date = new Date(today);
  }

  function smallCalendarClick(clickedDate: Date) {
    const range = getVisibleRange(date, view);
    if (isInRange(clickedDate, range.start, range.end) && clickedDate.getMonth() === date.getMonth()) {
      showDateModal(clickedDate, repository.events
        .filter((event) => event.date.start.getTime() <= clickedDate.getTime() + 24 * 60 * 60 * 1000 && event.date.end.getTime() >= clickedDate.getTime())
        .sort(compareEventsByStartDate)
      );
    } else {
      date = clickedDate;
    }
  }

  /* Fetching logic */
  let isLoading: boolean = $derived(getMetadata().loadingData);
  let loaderAnimation = $state(false);
  $effect(() => {
    if (isLoading) loaderAnimation = true;
  })

  afterNavigate(() => {
    refresh();
  });

  beforeNavigate((args) => {
    if (args.to === null) return;
    clearTimeout(spooledRefresh);
    spooledRefresh = undefined; 
  });

  let spooledRefresh: (ReturnType<typeof setTimeout> | undefined) = $state(undefined);
  function refresh(force = false) {
    const range = getVisibleRange(date, view);

    getRepository().getAllEvents(range.start, range.end, force).catch((err) => {
      queueNotification(ColorKeys.Danger, `Failed to fetch events: ${err.message}`);
    });

    connectivity.check();

    if (!publicReadonly) {
      clearTimeout(spooledRefresh);
      spooledRefresh = setTimeout(() => {
        refresh();
      }, autoRefreshInterval);
    }
  }

  function forceRefresh() {
    getRepository().invalidateCache();
    refresh(true);
  }

  $effect(() => {
    ((date: Date, view: "month" | "week" | "day" | "agenda") => {
      untrack(() => {
        if (!browser) return;
        refresh();
      });
    })(date, view);
  });

  onMount(() => {
    if (!browser) return;
    getRepository().getSources().catch(NoOp);
  });

  /* Single instance modal logic */
  let showSourceWizardModal: () => any = $state(NoOp);

  let showNewSourceModal: () => any = $state(NoOp);
  const showNewSourceModalInternal = () => { return showNewSourceModal(); };
  setContext("showNewSourceModal", showNewSourceModalInternal);

  let showSourceModal: (source: SourceModel) => any = $state(NoOp);
  const showSourceModalInternal = (source: SourceModel) => { return showSourceModal(source); };
  setContext("showSourceModal", showSourceModalInternal);

  let showNewCalendarModal: () => any = $state(NoOp);

  let showCalendarModal: (calendar: CalendarModel) => any = $state(NoOp);
  const showCalendarModalInternal = (calendar: CalendarModel) => { return showCalendarModal(calendar); };
  setContext("showCalendarModal", showCalendarModalInternal);

  let showNewEventModal: (date: Date) => any = $state(NoOp);
  const showNewEventModalInternal = (date: Date) => { return showNewEventModal(date); };
  setContext("showNewEventModal", showNewEventModalInternal);

  let showEventModal: (event: EventModel) => any = $state(NoOp);
  const showEventModalInternal = (event: EventModel) => { return showEventModal(event); };
  setContext("showEventModal", showEventModalInternal);

  let showDateModal: (date: Date, events: (EventModel | null)[]) => any = $state(NoOp);
  const showDateModalInternal = (date: Date, events: (EventModel | null)[]) => { return showDateModal(date, events); };
  setContext("showDateModal", showDateModalInternal);

  let showSettingsModal: () => any = $state(NoOp);

  let showCreditsModal: () => any = $state(NoOp);

  const timezone = $derived(Intl.DateTimeFormat().resolvedOptions().timeZone || "UTC");
  const connectivityLabel = $derived.by(() => {
    switch (connectivity.reachable) {
      case Reachability.Database:
        return "Sync OK";
      case Reachability.Backend:
        return "DB indisponible";
      case Reachability.Frontend:
        return "Backend indisponible";
      case Reachability.None:
        return "Frontend indisponible";
      case Reachability.Incompatible:
        return "Versions incompatibles";
      default:
        return "Etat inconnu";
    }
  });
  const nextEventLabel = $derived.by(() => {
    const now = Date.now();
    const next = repository.events
      .filter((event) => event.date.end.getTime() > now)
      .sort(compareEventsByStartDate)[0];

    if (!next) return "Prochain: aucun";
    const time = next.date.start.toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" });
    return `Prochain: ${time} ${next.name}`;
  });

  const agendaDays = $derived.by(() => {
    const range = getVisibleRange(date, view);
    const start = new Date(range.start);
    start.setHours(0, 0, 0, 0);
    const end = new Date(range.end);
    end.setHours(23, 59, 59, 999);

    const inRangeEvents = repository.events
      .filter((event) => event.date.end.getTime() > start.getTime() && event.date.start.getTime() <= end.getTime())
      .sort(compareEventsByStartDate);

    const byDay = new Map<string, { date: Date, events: EventModel[] }>();
    for (const event of inRangeEvents) {
      const day = new Date(event.date.start);
      day.setHours(0, 0, 0, 0);
      const key = day.toISOString().split("T")[0];
      if (!byDay.has(key)) byDay.set(key, { date: day, events: [] });
      byDay.get(key)?.events.push(event);
    }

    return Array.from(byDay.values()).sort((a, b) => a.date.getTime() - b.date.getTime());
  });

  const calendarGranularity = $derived(view === "agenda" ? "month" : view);
  const todayKey = $derived(new Date().toISOString().split("T")[0]);
  let agendaContainer: HTMLElement | null = $state(null);

  $effect(() => {
    if (view !== "agenda" || agendaDays.length === 0) return;
    setTimeout(() => {
      const todayRow = agendaContainer?.querySelector<HTMLElement>("[data-agenda-today='true']");
      todayRow?.scrollIntoView({ block: "start", behavior: "auto" });
    }, 0);
  });

  function openAgendaEvent(event: EventModel) {
    showEventModal(event).then((updatedEvent: EventModel) => event = updatedEvent).catch(NoOp);
  }

  function agendaItemKeydown(e: KeyboardEvent, event: EventModel) {
    if (e.key === "Enter" || e.key === " ") {
      e.preventDefault();
      openAgendaEvent(event);
    }
  }

  function agendaParticipantColors(event: EventModel): string[] {
    const colors = new Set(event.participant_colors || []);
    const own = GetEventColor(event);
    if (own) colors.add(own);
    return Array.from(colors);
  }
</script>

<style lang="scss">
  @use "../../styles/animations.scss";
  @use "../../styles/colors.scss";
  @use "../../styles/dimensions.scss";
  @use "../../styles/text.scss";

  :global(body) {
    display: flex;
    flex-direction: row;
    //display: grid;
    //grid-template-columns: auto 1fr;
    ////grid-template-rows: 1fr auto;
    ////grid-template-areas:
    ////  "aside main"
    ////  "aside footer";
    //grid-template-rows: auto;
    //grid-template-areas: "aside main";
  }

  main {
    width: 100%;
    height: 100%;
    display: flex;
    flex-direction: column;
    gap: dimensions.$gapSmall;
    grid-area: main;
    border: 1px solid var(--colorBorderSubtle, #3a3a3a);
    background-color: colors.$backgroundPrimary;
    overflow: hidden;
  }
  
  div.leftPane {
    display: flex;
    flex-direction: row;
    gap: 0;
  }

  nav.activityBar {
    width: 48px;
    min-width: 48px;
    max-width: 48px;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: dimensions.$gapSmall;
    background-color: colors.$backgroundTertiary;
    border: 1px solid var(--colorBorderSubtle, #3a3a3a);
    padding: dimensions.$gapSmall dimensions.$gapSmaller;
  }

  button.activityIcon {
    all: unset;
    width: 100%;
    height: 34px;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    color: colors.$foregroundSecondary;
    border-left: 2px solid transparent;
  }

  button.activityIcon:hover {
    background-color: var(--colorBackgroundHover, #2a2d2e);
    color: colors.$foregroundPrimary;
  }

  button.activityIcon.active {
    color: colors.$foregroundPrimary;
    border-left-color: colors.$backgroundAccent;
    background-color: var(--colorBackgroundSelection, #094771);
  }

  aside.sidebar {
    display: flex;
    flex-direction: column;
    gap: dimensions.$gapSmall;
    min-width: 15em;
    width: 15em;
    max-width: 15em;
    grid-area: aside;
    background-color: colors.$backgroundSecondary;
    border-top: 1px solid var(--colorBorderSubtle, #3a3a3a);
    border-right: 1px solid var(--colorBorderSubtle, #3a3a3a);
    border-bottom: 1px solid var(--colorBorderSubtle, #3a3a3a);
    border-left: 0;
    padding: dimensions.$gapSmall;
  }

  div.sources {
    flex-grow: 1;
    display: flex;
    flex-direction: column;
    gap: dimensions.$gapSmall;
    overflow: auto;
    margin: 0;
    padding: dimensions.$gapTiny;
  }

  div.toprow {
    display: flex;
    flex-direction: row;
    gap: dimensions.$gapSmall;
    justify-content: space-between;
    margin: 0;
    align-items: center;
    padding: dimensions.$gapSmall;
    border-bottom: 1px solid var(--colorBorderSubtle, #3a3a3a);
    background-color: colors.$backgroundTertiary;
  }

  span.reachability {
    color: colors.$backgroundFailure;
    align-items: center;
    display: flex;
    flex-direction: row;
    justify-content: center;
    gap: dimensions.$gapSmall;
  }

  span.refreshButtonWrapper {
    display: flex;
    align-items: center;
    justify-content: center;
  }

  span.spin {
    animation: spin animations.$animationSpeedSlow animations.$cubic infinite forwards;
  }

  span.copyright {
    color: color-mix(in srgb, colors.$foregroundPrimary 50%, transparent);
    font-size: text.$fontSizeSmall;
    text-align: center;
    margin-top: 0;
  }

  div.statusbar {
    min-height: 22px;
    height: 22px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: dimensions.$gapMiddle;
    padding: 0 dimensions.$gapSmall;
    background-color: colors.$backgroundAccent;
    color: colors.$foregroundAccent;
    font-size: text.$fontSizeSmall;
    border-top: 1px solid color-mix(in srgb, colors.$foregroundAccent 30%, transparent);
    white-space: nowrap;
    overflow: hidden;
  }

  div.statusbar span {
    overflow: hidden;
    text-overflow: ellipsis;
  }

  section.agenda {
    height: 100%;
    overflow-y: auto;
    padding: 0;
    display: flex;
    flex-direction: column;
    gap: 0;
  }

  div.agendaDay {
    border-bottom: 1px solid var(--colorBorderSubtle, #3a3a3a);
    background-color: colors.$backgroundSecondary;
  }

  h3.agendaDate {
    margin: 0;
    padding: dimensions.$gapSmall;
    font-size: text.$fontSize;
    font-weight: text.$fontWeightTitle;
    background-color: colors.$backgroundTertiary;
    border-bottom: 1px solid var(--colorBorderSubtle, #3a3a3a);
  }

  div.agendaItem {
    display: grid;
    grid-template-columns: 4.5em 1fr;
    gap: dimensions.$gapSmall;
    align-items: center;
    padding: dimensions.$gapSmaller dimensions.$gapSmall;
    border-left: 2px solid colors.$backgroundAccent;
    cursor: pointer;
  }

  div.agendaItem:hover {
    background-color: var(--colorBackgroundHover, #2a2d2e);
  }

  div.agendaItem:focus-visible {
    outline: 1px solid colors.$backgroundAccent;
    outline-offset: -1px;
  }

  span.agendaTime {
    color: colors.$foregroundSecondary;
    font-family: text.$fontFamilyTime;
    font-size: text.$fontSize;
    font-weight: text.$fontWeightSemiBold;
  }

  span.agendaName {
    display: inline-flex;
    align-items: center;
    gap: dimensions.$gapSmall;
  }

  span.eventDot {
    display: inline-block;
    width: 0.55em;
    height: 0.55em;
    border-radius: 50%;
    flex-shrink: 0;
  }

</style>

{#if !publicReadonly}
  <SourceWizardModal bind:showModal={showSourceWizardModal}/>
  <SourceModal bind:showCreateModal={showNewSourceModal} bind:showModal={showSourceModal}/>
  <CalendarModal bind:showCreateModal={showNewCalendarModal} bind:showModal={showCalendarModal}/>
  <DayViewModal bind:showModal={showDateModal}/>
{/if}
<EventModal bind:showCreateModal={showNewEventModal} bind:showModal={showEventModal}/>
<SettingsModal bind:showModal={showSettingsModal} appearanceOnly={publicReadonly}/>
<CreditsModal bind:showModal={showCreditsModal}/>

<div class="leftPane">
  <nav class="activityBar" aria-label="Vues calendrier">
    <button class="activityIcon" class:active={view === "month"} onclick={() => view = "month"} title="Mois">
      <CalendarDays size={16}/>
    </button>
    <button class="activityIcon" class:active={view === "week"} onclick={() => view = "week"} title="Semaine">
      <CalendarRange size={16}/>
    </button>
    <button class="activityIcon" class:active={view === "day"} onclick={() => view = "day"} title="Jour">
      <CalendarSingle size={16}/>
    </button>
    <button class="activityIcon" class:active={view === "agenda"} onclick={() => { date = new Date(); view = "agenda"; }} title="Agenda">
      <List size={16}/>
    </button>
  </nav>

  <aside class="sidebar">
    {#if settings.userSettings[UserSettingKeys.DisplaySmallCalendar]}
      <SmallCalendar date={date} smaller={true} onDayClick={(clickedDate) => smallCalendarClick(clickedDate)}></SmallCalendar>
    {/if}

    <div class="sources">
      {@render sourceEntries(repository.sources)}
    </div>

    <Horizontal position="center">
      <IconButton click={showSettingsModal}>
        <Settings/>
      </IconButton>
      {#if !publicReadonly}
        <IconButton click={showSourceWizardModal}>
          <PlusIcon/>
        </IconButton>
        <IconButton click={showCreditsModal}>
          <Copyleft/>
        </IconButton>
      {/if}
    </Horizontal>

    <span class="copyright">
      Copyright © 2026 Kacper Darowski (Opisek)<br>
      Licensed under TBD 
    </span>
  </aside>
</div>

<main>
  <div class="toprow">
    <MonthSelection bind:date granularity={calendarGranularity} />
    <Horizontal position="justify" width="auto">
      {#if connectivity.reachable != Reachability.Database}
        <span class="reachability">
          {#if connectivity.reachable == Reachability.Backend}
            The database cannot be reached.
          {:else if connectivity.reachable == Reachability.Frontend}
            The backend server cannot be reached.
          {:else if connectivity.reachable == Reachability.None}
            The frontend server cannot be reached.
          {:else if connectivity.reachable == Reachability.Incompatible}
            The frontend server and the backend server are not compatible.
          {:else}
            Unknown network error
          {/if}
          <WifiOff size={20}/>
        </span>
      {/if}

      <IconButton click={forceRefresh}>
        <span class="refreshButtonWrapper" class:spin={loaderAnimation} onanimationiteration={() => { if (!isLoading) loaderAnimation = false; }}>
          <RefreshCw size={20}/>
        </span>
      </IconButton>

      {#if !todayInRange}
        <Button onClick={seeToday} compact={true}>
          Today
        </Button>
      {/if}

      {#if !settings.userSettings[UserSettingKeys.ThemeSynchronize]}
        <ThemeToggle/>
      {/if}

    </Horizontal>
  </div>
  {#if view === "agenda"}
    <section class="agenda" bind:this={agendaContainer}>
      {#if agendaDays.length === 0}
        <div class="agendaDay">
          <h3 class="agendaDate">Aucun rendez-vous sur cette periode</h3>
        </div>
      {:else}
        {#each agendaDays as day (day.date.getTime())}
          <div
            class="agendaDay"
            data-agenda-today={day.date.toISOString().split("T")[0] === todayKey ? "true" : "false"}
          >
            <h3 class="agendaDate">{day.date.toLocaleDateString([], { weekday: "long", day: "2-digit", month: "long", year: "numeric" })}</h3>
            {#each day.events as event (event.id + event.date.start.getTime())}
              <div
                class="agendaItem"
                role="button"
                tabindex="0"
                onclick={() => openAgendaEvent(event)}
                onkeydown={(e) => agendaItemKeydown(e, event)}
              >
                <span class="agendaTime">
                  {#if event.date.allDay}
                    Journee
                  {:else}
                    {event.date.start.toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" })}
                  {/if}
                </span>
                <span class="agendaName">
                  {#each agendaParticipantColors(event) as color (color)}
                    <span class="eventDot" style="background-color: {color}"></span>
                  {/each}
                  {event.name}
                </span>
              </div>
            {/each}
          </div>
        {/each}
      {/if}
    </section>
  {:else}
    <Calendar
      date={date}
      view={view}
      events={repository.events}
      readOnly={publicReadonly}
    />
  {/if}
  <div class="statusbar">
    <span>{timezone}</span>
    <span>{connectivityLabel}</span>
    <span>{nextEventLabel}</span>
  </div>
</main>

{#snippet sourceEntries(sources: SourceModel[])}
  {#each sources as source, i (source.id)}
    <SourceEntry bind:source={repository.sources[i]} readOnly={publicReadonly}/>
    {#if !metadata.collapsedSources.has(repository.sources[i].id)}
      {@render calendarEntries(repository.calendars.filter(cal => cal.source === source.id) || [])}
    {/if}
  {/each}
{/snippet}

{#snippet calendarEntries(calendars: CalendarModel[])}
  {#each calendars as cal (cal.id)}
    {@const index = repository.calendars.findIndex((calendar) => calendar.id === cal.id)}
    <CalendarEntry bind:calendar={repository.calendars[index]} readOnly={publicReadonly}/>
  {/each}
{/snippet}