<script lang="ts">
  import type { Operation } from "./ui-generator";
  import api from "./api";

  import UIGen from "./UIGenerator.svelte";

  let selectedGroupOp: Operation[];
  let selectedOp: Operation;
</script>

<p>
  In order to work with the API you have to set the authentication token in the
  input box before executing any operation
</p>
<p>Token: <input bind:value={api.authToken} type="password" size="48" /></p>

<p>
  Operation:
  <select bind:value={selectedGroupOp}>
    <option selected />
    {#each Object.keys(api.operations) as group}
      <option value={api.operations[group]}>{group}</option>
    {/each}
  </select>
  {#if selectedGroupOp}
    <select bind:value={selectedOp}>
      <option selected />
      {#each selectedGroupOp as op}
        <option value={op}>{op.name}</option>
      {/each}
    </select>
  {/if}
</p>
<hr />
<p>
  {#if selectedOp}
    <UIGen operation={selectedOp} />
  {/if}
</p>
