import discord
import os

from discord import app_commands
from dotenv import load_dotenv
from typing import Optional

# =================================================================================================================================================
# Initial setup
# =================================================================================================================================================

intents = discord.Intents.default()
bot = discord.Client(intents = intents)
tree = discord.app_commands.CommandTree(bot)

# =================================================================================================================================================
# Slash Commands
# =================================================================================================================================================

# Tasks - Top Level Group
# =================================================================================================================================================
tasks = app_commands.Group(
    name="tasks",
    description="Everything related to tasks."
)

# Tasks/Show - Sub-command
# =================================================================================================================================================
show = app_commands.Group(
    name="show",
    description="Show tasks",
    parent=tasks
)

# Tasks/Show/All - Annotations
# =================================================================================================================================================
@show.command(name = "all", description = "Show all tasks.")
# Tasks/Show/All - Function Definition
# =================================================================================================================================================
async def tasks_show_all(interaction: discord.Interaction):
    await interaction.response.send_message("Showing all tasks...")

# Tasks/Show/Trader - Annotations
# =================================================================================================================================================
@show.command(name = "trader", description = "Show all tasks given by a trader.")
@app_commands.choices(
    trader = [
        app_commands.Choice(name = "Prapor", value = "prapor"),
        app_commands.Choice(name = "Therapist", value = "therapist"),
        app_commands.Choice(name = "Fence", value = "fence"),
        app_commands.Choice(name = "Skier", value = "skier"),
        app_commands.Choice(name = "Peacekeeper", value = "peacekeeper"),
        app_commands.Choice(name = "Mechanic", value = "mechanic"),
        app_commands.Choice(name = "Ragman", value = "ragman"),
        app_commands.Choice(name = "Jaeger", value = "jaeger"),
        app_commands.Choice(name = "BTR Driver", value = "btr_driver"),
    ]
)
# Tasks/Show/Trader - Function Definition
# =================================================================================================================================================
async def tasks_show_trader(interaction: discord.Interaction, trader: app_commands.Choice[str]):
    await interaction.response.send_message(f"Showing all tasks given by a {trader}...")

# =================================================================================================================================================
# Run Bot
# =================================================================================================================================================

# add our commands to the tree
tree.add_command(tasks)

# load .env file
load_dotenv()

# ready the bot
@bot.event
async def on_ready():
    print(f"Logged in as {bot.user} (ID: {bot.user.id})")
    print("Syncing commands...")
    await tree.sync()
    print("Commands synced in all servers!")

    # change bot activity
    await bot.change_presence(status = discord.Status.online, activity = discord.Game("Hiding from USECs"))

# run the bot
bot.run(os.getenv("TOKEN"))