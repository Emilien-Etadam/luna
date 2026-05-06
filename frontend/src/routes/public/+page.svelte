<script lang="ts">
  import { onMount } from "svelte";
  import { browser } from "$app/environment";

  let calendarEl = $state<HTMLDivElement | null>(null);
  let publicationDisabled = $state<boolean | null>(null);

  onMount(() => {
    if (!browser) {
      return () => {};
    }

    let destroy: (() => void) | null = null;

    void (async () => {
      const probe = await fetch("/api/public/all.ics", { method: "HEAD" });
      if (probe.status === 404) {
        publicationDisabled = true;
        return;
      }
      publicationDisabled = false;

      if (!calendarEl) {
        return;
      }

      const [core, dayGrid, timeGrid, listPlugin, iCalendar, frLocale] = await Promise.all([
        import("@fullcalendar/core"),
        import("@fullcalendar/daygrid"),
        import("@fullcalendar/timegrid"),
        import("@fullcalendar/list"),
        import("@fullcalendar/icalendar"),
        import("@fullcalendar/core/locales/fr"),
      ]);

      const CalendarCtor = core.Calendar;
      const calendar = new CalendarCtor(calendarEl, {
        plugins: [dayGrid.default, timeGrid.default, listPlugin.default, iCalendar.default],
        locale: frLocale.default,
        firstDay: 1,
        initialView: "dayGridMonth",
        headerToolbar: {
          left: "prev,next today",
          center: "title",
          right: "dayGridMonth,timeGridWeek,listWeek",
        },
        events: {
          url: "/api/public/all.ics",
          format: "ics",
        },
      });

      calendar.render();
      destroy = () => calendar.destroy();
    })();

    return () => {
      destroy?.();
    };
  });
</script>

{#if publicationDisabled === true}
  <p class="msg">Publication désactivée</p>
{:else if publicationDisabled === null}
  <p class="msg">Chargement…</p>
{/if}

{#if publicationDisabled !== true}
  <div bind:this={calendarEl} class="fc-root"></div>
{/if}

<style>
  .fc-root {
    min-height: 70vh;
    margin: 1rem;
  }

  .msg {
    margin: 2rem;
    font-family:
      system-ui,
      -apple-system,
      sans-serif;
  }
</style>
