 <script lang="ts">
 	import { journalApi } from '$lib/api/journal';
 	import type { CreateEntryRequest } from '$lib/types/journal';

	let title = '';
	let content = '';
	let saving = false;
	let error: string | null = null;

	async function createEntry() {
		if (!title.trim() || !content.trim()) {
			alert('Please fill in both title and content');
			return;
		}

		saving = true;
		error = null;

		try {
			const entryData: CreateEntryRequest = {
				title: title.trim(),
				content: content.trim()
			};
			const newEntry = await journalApi.createEntry(entryData);
			window.location.href = `/entry/${newEntry.id}`;
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to create entry';
		} finally {
			saving = false;
		}
	}

	function cancel() {
		if (title.trim() || content.trim()) {
			if (confirm('Are you sure you want to discard your entry?')) {
				window.location.href = '/';
			}
		} else {
			window.location.href = '/';
		}
	}
</script>

<div class="max-w-4xl mx-auto">
 	<div class="bg-white border rounded-lg shadow">
 		<div class="bg-gradient-to-r from-green-500 to-emerald-600 text-white rounded-t-lg p-6">
 			<h1 class="text-3xl font-bold">New Journal Entry</h1>
 			<p class="text-green-100 mt-2">Capture your thoughts and memories</p>
 		</div>
 		<div class="p-8">
 			{#if error}
 				<div class="bg-red-50 border border-red-200 rounded-lg p-4 mb-6">
 					<p class="text-red-800">{error}</p>
 				</div>
 			{/if}

 			<form onsubmit={createEntry} class="space-y-6">
 				<div class="space-y-2">
 					<label for="title" class="text-lg font-semibold text-gray-800">Title</label>
 					<input
 						id="title"
 						type="text"
 						bind:value={title}
 						placeholder="Give your entry a title..."
 						required
 						class="w-full px-3 py-2 border border-gray-300 rounded-md text-lg"
 					/>
 				</div>

 				<div class="space-y-2">
 					<label for="content" class="text-lg font-semibold text-gray-800">Content</label>
 					<textarea
 						id="content"
 						bind:value={content}
 						rows={12}
 						placeholder="Share your thoughts, feelings, and experiences..."
 						required
 						class="w-full px-3 py-2 border border-gray-300 rounded-md text-lg leading-relaxed resize-none"
 					></textarea>
 				</div>

 				<div class="flex justify-between items-center pt-4">
 					<button type="button" onclick={cancel} class="px-6 py-3 border border-gray-300 rounded-lg hover:bg-gray-50">
 						Cancel
 					</button>
 					<button type="submit" disabled={saving} class="px-8 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50">
 						{#if saving}
 							Creating...
 						{:else}
 							Create Entry
 						{/if}
 					</button>
 				</div>
 			</form>
 		</div>
 	</div>
 </div>