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

# Tasks - Annotations
# =================================================================================================================================================
@tree.command(name = "tasks", description = "Get information regarding tasks required for Kappa.")
@app_commands.describe(given_by = "Choose to show tasks by individual trader.")
@app_commands.choices(
    given_by = [
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
# Tasks - Function Definition
# =================================================================================================================================================
async def tasks(interaction: discord.Interaction, given_by: Optional[app_commands.Choice[str]]):

    if given_by != None:

        await interaction.response.send_message(
            f"Selected trader: {given_by.name}"
        )
    
    else:
        
        await interaction.response.send_message(
            f"All"
        )

# =================================================================================================================================================
# Run Bot
# =================================================================================================================================================

# load .env file
load_dotenv()

# ready the bot
@bot.event
async def on_ready():
    print(f"Logged in as {bot.user} (ID: {bot.user.id})")
    print("Syncing commands...")
    await tree.sync(guild = discord.Object(id = os.getenv("GUILD_ID")))
    print("Commands synced in all servers!")

    # change bot activity
    await bot.change_presence(status = discord.Status.online, activity = discord.Game("Hiding from USECs"))

# run the bot
bot.run(os.getenv("TOKEN"))