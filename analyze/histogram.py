import sys
from PIL import Image
import matplotlib.pyplot as plt

def plot_histograms(original_path: str, stego_path: str):
    # Загружаем изображения
    orig = Image.open(original_path).convert("RGBA")
    stego = Image.open(stego_path).convert("RGBA")

    # Получаем каналы
    orig_r, orig_g, orig_b, _ = orig.split()
    stego_r, stego_g, stego_b, _ = stego.split()

    # Строим гистограммы
    fig, axs = plt.subplots(3, 2, figsize=(12, 9))
    axs = axs.flatten()

    # Канал R
    axs[0].hist(list(orig_r.getdata()), bins=256, color='red', alpha=0.7)
    axs[0].set_title("R - до")
    axs[1].hist(list(stego_r.getdata()), bins=256, color='red', alpha=0.7)
    axs[1].set_title("R - после")

    # Канал G
    axs[2].hist(list(orig_g.getdata()), bins=256, color='green', alpha=0.7)
    axs[2].set_title("G - до")
    axs[3].hist(list(stego_g.getdata()), bins=256, color='green', alpha=0.7)
    axs[3].set_title("G - после")

    # Канал B
    axs[4].hist(list(orig_b.getdata()), bins=256, color='blue', alpha=0.7)
    axs[4].set_title("B - до")
    axs[5].hist(list(stego_b.getdata()), bins=256, color='blue', alpha=0.7)
    axs[5].set_title("B - после")

    plt.tight_layout()
    plt.savefig("histogram.png")

if __name__ == "__main__":
    plot_histograms(sys.argv[1], sys.argv[2])
