<script lang="ts">
	var websocket = new WebSocket(window.location.origin+'/socket');
	let todoItems = $state({});
	let newToDoName = $state('');
	let newToDoDescription = $state('');

	websocket.onmessage = function (e) {
		let data = JSON.parse(e.data);
		if (data.type == 'clear') {
			todoItems = {};
		} else if (data.type == 'remove') {
			todoItems = Object.fromEntries(
				Object.entries(todoItems).filter(([key, value]) => key != data.key)
			);
		} else {
			let data = JSON.parse(e.data);
			todoItems = Object.assign({}, todoItems, data);
		}
	};

	function remove(key: string) {
		let data = {
			type: 'remove',
			key: key
		};
		websocket.send(JSON.stringify(data));
	}

	function add() {
		let data = {
			type: 'add',
			name: newToDoName,
			description: newToDoDescription
		};
		websocket.send(JSON.stringify(data));
		newToDoName = '';
		newToDoDescription = '';
	}

	function done(key: string) {
		let data = {
			type: 'done',
			key: key
		};
		websocket.send(JSON.stringify(data));
	}
</script>

<div>
	<div>
		<h2>Websocket Todo App</h2>
	</div>
	<div class="input-area">
		<div class="input-area-name">
			<div>Name</div>
			<input class="input-field-name" bind:value={newToDoName} />
		</div>
		<div class="input-area-description">
			<div>Description</div>
			<input class="input-field-description" bind:value={newToDoDescription} />
		</div>
		<button class="add-btn" onclick={add}>
			<img src="add.svg" alt="Add Icon" />
		</button>
	</div>
	<div class="main_area">
		{#each Object.entries(todoItems) as [key, value]}
			<div class="item-dev {(value as any).done ? 'item-dev-done' : 'item-dev-todo'}">
				{#if !(value as any).done}
					<button class="done-btn" onclick={() => done(key)}>
						<img src="check.svg" alt="Check Icon" />
					</button>
					<div class="devider"></div>
				{/if}

				<div>
					<div class="name">
						{(value as any).name}
					</div>
					<div>
						{(value as any).description}
					</div>
				</div>
				<div class="date">
					{(value as any).date}
				</div>

				<button class="remove-btn" onclick={() => remove(key)}>
					<img src="delete.svg" alt="Delete Icon" />
				</button>
			</div>
		{/each}
	</div>
</div>

<style>
	.input-field-description:focus-visible {
		outline: 1px solid cyan;
	}
	.input-field-name:focus-visible {
		outline: 1px solid cyan;
	}
	.main_area {
		display: flex;
		flex-direction: column;
	}
	.input-field-name {
		background: #374151;
		border: none;
		color: cyan;
		height: 3vh;
	}

	.input-field-description {
		background: #374151;
		border: none;
		color: cyan;
		height: 3vh;
	}
	.input-area {
		display: flex;
		margin: 5px;
		align-items: end;
	}
	.input-area-name {
		display: flex;
		flex-direction: column;
		width: 20%;
	}
	.input-area-description {
		display: flex;
		flex-direction: column;
		width: -webkit-fill-available;
		margin: 0px 5px;
	}
	.name {
		font-size: 21px;
	}
	.devider {
		height: 100%;
		width: 2px;
		background: #22d3ee;
		margin: 0px 10px;
	}
	.item-dev {
		display: flex;
		border-radius: 5px;
		padding: 10px;
		margin: 5px;
		height: 4vh;
	}

	.item-dev-done {
		border: 1px solid #39ff14;
	}

	.item-dev-todo {
		border: 1px solid cyan;
	}

	.remove-btn {
		margin-left: auto;
		background: transparent;
		border: 2px solid cyan;
		border-radius: 5px;
		width: 4vh;
	}

	.done-btn {
		background: transparent;
		border: 2px solid cyan;
		border-radius: 5px;
		width: 4vh;
	}

	.add-btn {
		background: transparent;
		border: 2px solid cyan;
		border-radius: 5px;
		width: 5vh;
		height: 4vh;
	}

	.date {
		font-size: 12px;
	}
</style>
