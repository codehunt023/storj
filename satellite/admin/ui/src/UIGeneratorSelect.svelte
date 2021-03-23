<script lang="ts">
  import type { Select } from "./ui-generator";

  export let label: string;
  export let config: Select;
  export let value: string | number;
</script>

<label>
  {label}
  {#if config.required}<sup>*</sup>{/if}:
  {#if config.multiple}
    <select bind:value required={config.required} multiple="true">
      {#each config.options as option}
        <option value={option.value}>{option.text}</option>
      {/each}
    </select>
  {:else}
    <select bind:value required={config.required}>
      {#each config.options as option}
        <option value={option.value}>{option.text}</option>
      {/each}
    </select>
  {/if}
</label>
