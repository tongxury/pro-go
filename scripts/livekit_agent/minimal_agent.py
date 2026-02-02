import asyncio
import os
from dotenv import load_dotenv
from livekit.agents import JobContext, JobRequest, WorkerOptions, cli, llm, voice_assistant
from livekit.plugins import openai, cartesia, silero

load_dotenv()

# Check environment variables
print(f"--- Environment Status ---")
print(f"LIVEKIT_URL: {'Set' if os.getenv('LIVEKIT_URL') else 'MISSING'}")
print(f"OPENAI_API_KEY: {'Set' if os.getenv('OPENAI_API_KEY') else 'MISSING'}")
print(f"CARTESIA_API_KEY: {'Set' if os.getenv('CARTESIA_API_KEY') else 'MISSING'}")
print(f"--------------------------")

async def request_fnc(req: JobRequest):
    print(f"--- Received job request for room: {req.job.room.name} ---")
    await req.accept()

async def entrypoint(ctx: JobContext):
    print(f"!!! Entrypoint triggered for room: {ctx.room.name} !!!")
    print(f"Room name: {ctx.room.name}")
    initial_ctx = llm.ChatContext().append(
        role="system",
        text="You are AURA, a supportive and professional AI mental health counselor. Keep your responses empathetic and concise.",
    )
    
    print(f"Connecting to room {ctx.room.name}...")
    await ctx.connect()
    print("Connected to room!")

    # Configure the voice assistant
    assistant = voice_assistant.VoiceAssistant(
        vad=silero.VAD.load(),
        stt=openai.STT(),
        llm=openai.LLM(),
        tts=cartesia.TTS(),
        chat_ctx=initial_ctx,
    )

    @assistant.on("tts_stream_started")
    def on_tts_start(stream):
        print(">>> TTS Stream Started!")

    @assistant.on("agent_speech_started")
    def on_speech_start():
        print(">>> AURA started speaking (Agent track publishing...)")

    @assistant.on("agent_speech_stopped")
    def on_speech_stop():
        print(">>> AURA stopped speaking.")

    @ctx.room.on("track_published")
    def on_track_published(publication, participant):
        print(f">>> Track published by {participant.identity}: {publication.sid} ({publication.source})")

    @ctx.room.on("participant_connected")
    def on_participant_connected(participant):
        print(f">>> Participant connected: {participant.identity}")

    assistant.start(ctx.room)
    
    # Wait for a brief moment to ensure track publication
    await asyncio.sleep(2)
    print(f"Assistant active. Local participant: {ctx.room.local_participant.identity}")
    
    # Safely check for microphone track
    is_publishing = any(pub.source == 'microphone' for pub in ctx.room.local_participant.track_publications.values())
    print(f"Is publishing audio: {is_publishing}")

    print("Assistant started! Greet user...")
    
    # Optional: Greet the user when they join
    await assistant.say("Hi, I'm AURA. I'm here to listen. How are you feeling today?", allow_interruptions=True)
    print("Greeting sent.")

if __name__ == "__main__":
    cli.run_app(WorkerOptions(
        entrypoint_fnc=entrypoint, 
        request_fnc=request_fnc, 
        port=8082,
        load_threshold=0.99 # Ignore high CPU load during development
    ))
    
# To run this:
# 1. pip install -r requirements.txt
# 2. export LIVEKIT_URL=...
# 3. export LIVEKIT_API_KEY=...
# 4. export LIVEKIT_API_SECRET=...
# 5. export OPENAI_API_KEY=...
# 6. export CARTESIA_API_KEY=...
# 7. python minimal_agent.py dev
