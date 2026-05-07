<script lang="ts">
  import type { Snippet } from "svelte";

  import Button from "../interactive/Button.svelte";
  import CloseButton from "../interactive/CloseButton.svelte";
  import Horizontal from "../layout/Horizontal.svelte";
  import Title from "../layout/Title.svelte";

  import { NoOp } from "$lib/client/placeholders";
  import { redrawNotifications } from "$lib/client/notifications";
  import { isChildOfModal } from "../../lib/common/misc";

  interface Props {
    title: string;
    focusElement?: HTMLElement | null;
    onModalHide?: any;
    onModalSubmit?: any;
    showModal?: () => any;
    hideModal?: () => any;
    resetFocus?: () => any;
    children?: Snippet;
    buttons?: Snippet;
    topButtons?: Snippet;
  }

  let {
    title,
    focusElement = null,
    onModalHide = NoOp,
    showModal = $bindable(),
    hideModal = $bindable(NoOp),
    resetFocus = $bindable(),
    onModalSubmit = hideModal,
    children,
    buttons,
    topButtons,
  }: Props = $props();

  let dialog: HTMLDialogElement;

  let visible = $state(false);

  let ignoreClickOutside = $state(false);
  function mouseDown(event: MouseEvent) {
    ignoreClickOutside = isChildOfModal(event.target as HTMLElement) && event.target !== dialog;
  }

  function clickOutside(event: MouseEvent) {
    if (!dialog || ignoreClickOutside) return;
    if (event.target === dialog) {
      hideModal();
      event.stopPropagation();
    }
  }

  $effect(() => {
    if (visible) dialog.showModal();
  });

  resetFocus = () => {
    if (focusElement) focusElement.focus();
    else dialog.focus();
  }
  showModal = () => {
    window.addEventListener("mousedown", mouseDown);
    window.addEventListener("click", clickOutside);
    visible = true
    setTimeout(resetFocus, 0);
    setTimeout(redrawNotifications, 0); // hacky way to make sure that notifications are always on the very top. sometimes has a visible blink. should revisit one day.
  }
  hideModal = () => {
    window.removeEventListener("mousedown", mouseDown);
    window.removeEventListener("click", clickOutside);
    dialog.close();
  }

  function submitInternal(event: Event) {
    event.preventDefault();
    onModalSubmit();
    return false;
  }
</script>

<style lang="scss">
  @use "../../styles/animations.scss";
  @use "../../styles/colors.scss";
  @use "../../styles/dimensions.scss";

  dialog {
    border: 1px solid var(--border-default);
    border-radius: var(--borderRadiusLarge);
    padding: 0;
    background-color: var(--surface-raised);
    color: var(--fg-primary);
    box-shadow: var(--shadow-3);
    max-width: min(90vw, 640px);
    max-height: 90vh;
  }
  :global(html[data-frost="true"]) dialog {
    background-color: color-mix(in srgb, var(--surface-raised) 80%, transparent);
    backdrop-filter: blur(var(--overlay-blur));
  }
  dialog::backdrop {
    background-color: var(--overlay-bg);
    backdrop-filter: blur(var(--overlay-blur));
  }

  dialog[open] {
    animation: zoom animations.$animationSpeed animations.$cubic forwards;
  }

  dialog:focus {
    outline: none;
  }
  
  form {
    padding: 20px 24px 24px 24px;
    border-radius: var(--borderRadiusLarge);
    display: flex;
    flex-direction: column;
    flex-wrap: nowrap;
    gap: dimensions.$gapMiddle;
    box-sizing: content-box;
    min-width: 30em;
    width: fit-content;
  }

  @media (max-width: 720px) {
    dialog {
      max-width: 100vw;
      max-height: 100vh;
      width: 100vw;
      height: 100vh;
      border-radius: 0;
      border: 0;
    }
    form {
      min-width: 0;
      width: 100%;
      padding: 16px;
    }
  }
</style>

<dialog
  bind:this={dialog}
  onclose={() => {visible = false; onModalHide();}}
  class:closed={visible}
>
  {#if visible}
    <form onsubmit={submitInternal}>
      <Horizontal>
        <Title>
          {title}
        </Title>
        {#if topButtons}
          <Horizontal position="right">
            {@render topButtons()}
            <CloseButton onClick={hideModal} />
          </Horizontal>
        {:else}
          <CloseButton onClick={hideModal} />
        {/if}
      </Horizontal>
      {@render children?.()}
      <Horizontal position="right">
        {#if buttons}{@render buttons()}{:else}
          <Button onClick={hideModal}>Close</Button>
        {/if}
      </Horizontal>
    </form>
  {/if}
</dialog>