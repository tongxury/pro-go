import json
import argparse
import os

# ==============================================================================
# Script: sora_prompt_gen.py
# Description: Generates high-quality Sora-2 visual prompts from product data.
# Usage: python3 sora_prompt_gen.py --product "Smart Coffee Maker" --features "Minimalist, Stainless Steel, App control"
# ==============================================================================

SORA_TEMPLATES = {
    "cinematic": "Cinematic 8k video, hyper-realistic, shallow depth of field, {product} in a {environment}. The lighting is {lighting}, emphasizing its {feature}. Shot on 35mm lens, slow motion.",
    "abstract": "Abstract close-up of {product}, focusing on {feature} as it magically {action} in a zero-gravity environment. Surrounding textures are {lighting} and iridescent, 8k resolution, Unreal Engine 5 render style.",
    "lifestyle": "A high-end lifestyle shot showing a person enjoying {product} in a {environment}. The atmosphere is {lighting}, soft focus in the background, professional color grading, ultra-sharp 4k Sora generation.",
    "macro": "Macro photography video of {product}, revealing the intricate texture of {feature}. Particles of {lighting} dust float in the air, creating a dreamlike bokeh effect."
}

def generate_prompts(product, features, environment="futuristic clean studio", lighting="warm golden hour"):
    feature_list = [f.strip() for f in features.split(",")]
    primary_feature = feature_list[0] if feature_list else "design"
    
    results = {}
    for style, template in SORA_TEMPLATES.items():
        # Heuristic actions for abstract style
        action = "assembles itself" if "tech" in product.lower() or "minimalist" in features.lower() else "floats elegantly"
        
        prompt = template.format(
            product=product,
            feature=primary_feature,
            environment=environment,
            lighting=lighting,
            action=action
        )
        results[style] = prompt
    return results

def main():
    parser = argparse.ArgumentParser(description="Generate Sora-2 Prompts for Ecommerce")
    parser.add_argument("--product", required=True, help="Product name")
    parser.add_argument("--features", required=True, help="Comma-separated product features")
    parser.add_argument("--env", default="minimalist zen office", help="Visual environment")
    parser.add_argument("--lighting", default="cinematic soft light", help="Lighting style")
    parser.add_argument("--json", action="store_true", help="Output as JSON")

    args = parser.parse_args()

    prompts = generate_prompts(args.product, args.features, args.env, args.lighting)

    if args.json:
        print(json.dumps(prompts, indent=4, ensure_ascii=False))
    else:
        print(f"\n--- Sora-2 Prompts for: {args.product} ---\n")
        for style, p in prompts.items():
            print(f"[{style.upper()}]\n{p}\n")

if __name__ == "__main__":
    main()
