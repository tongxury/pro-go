import json
import argparse
import os

# ==============================================================================
# Script: sora_prompt_gen.py
# Description: Generates high-quality Sora-2 prompts with consistent IP profiles.
# ==============================================================================

SORA_TEMPLATES = {
    "cinematic": "Located in {environment}, cinematic 8k video, hyper-realistic, {product}. The model has {hair} and is wearing {accessories}. Visible {tattoo}. She is seen {signature_action}. The lighting is {lighting}, mood is {vibe}. Shot on 35mm lens. {ui_layout}",
    "lifestyle": "In the setting of {environment}, a high-end lifestyle shot showing a person with {hair} and wearing {accessories} enjoying {product}. They have a {tattoo} and are {signature_action}. Atmosphere is {lighting}, {vibe}. {ui_layout}",
    "fashion_consistent": "Showcased in {environment}, fashion show style: A female model with {hair} and {accessories} wearing {product}. She has a {tattoo} and pauses to {signature_action}. Lighting is {lighting}. {ui_layout}",
    "macro_detail": "Against the background of {environment}, macro photography video of {product}, revealing textures of {feature}. In the background, the blurred silhouette of the person with {hair}, {accessories}, and {tattoo} performing {signature_action} is visible. {ui_layout}"
}

def load_profiles(path):
    if os.path.exists(path):
        with open(path, 'r', encoding='utf-8') as f:
            return json.load(f)
    return {}

def generate_prompts(product, features, profile, scene_action="walking quietly"):
    feature_list = [f.strip() for f in features.split(",")]
    primary_feature = feature_list[0] if feature_list else "design"
    
    results = {}
    for style, template in SORA_TEMPLATES.items():
        prompt = template.format(
            product=product,
            feature=primary_feature,
            hair=profile.get("hair", ""),
            accessories=profile.get("accessories", ""),
            tattoo=profile.get("tattoo", ""),
            signature_action=profile.get("signature_action", ""),
            ui_layout=profile.get("ui_layout", ""),
            lighting=profile.get("lighting", ""),
            vibe=profile.get("vibe", ""),
            environment=profile.get("fixed_environment", "minimalist studio"),
            scene_action=scene_action
        )
        results[style] = prompt
    return results

def main():
    profile_path = os.path.join(os.path.dirname(__file__), "ip_profiles.json")
    profiles = load_profiles(profile_path)
    
    parser = argparse.ArgumentParser(description="Generate Sora-2 Prompts with IP Consistency Profiles")
    parser.add_argument("--product", required=True, help="Product name")
    parser.add_argument("--features", required=True, help="Comma-separated product features")
    parser.add_argument("--profile", default="default_ip", help=f"IP profile to use ({', '.join(profiles.keys())})")
    parser.add_argument("--env", default="minimalist zen office", help="Visual environment")
    parser.add_argument("--json", action="store_true", help="Output as JSON")

    args = parser.parse_args()

    if args.profile not in profiles:
        print(f"Warning: Profile '{args.profile}' not found. Using default.")
        profile = profiles.get("default_ip", {})
    else:
        profile = profiles[args.profile]

    prompts = generate_prompts(args.product, args.features, profile, args.env)

    if args.json:
        print(json.dumps(prompts, indent=4, ensure_ascii=False))
    else:
        print(f"\n--- Sora-2 Consistent IP Prompts [{args.profile.upper()}] for: {args.product} ---\n")
        for style, p in prompts.items():
            print(f"[{style.upper()}]\n{p}\n")

if __name__ == "__main__":
    main()
