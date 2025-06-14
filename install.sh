#!/bin/sh

set -e

echo "🚀 aperfetch kurulumu başlıyor..."

echo "🔍 Rust araç zinciri kontrol ediliyor..."
if ! command -v cargo > /dev/null; then
    echo "❌ Hata: Rust ve Cargo yüklü değil."
    echo "Lütfen önce https://rustup.rs/ adresinden Rust'ı kurun ve tekrar deneyin."
    echo "Termux için 'pkg install rust' komutunu kullanabilirsiniz."
    exit 1
fi
echo "✅ Rust araç zinciri bulundu."

echo "📦 aperfetch derleniyor... (Bu işlem biraz zaman alabilir)"
cargo build --release
echo "✅ Derleme tamamlandı."

INSTALL_DIR=""

if [ -n "$TERMUX_VERSION" ]; then
    echo "✅ Termux ortamı algılandı."
    INSTALL_DIR="$PREFIX/bin"
else
    INSTALL_DIR="$HOME/.local/bin"
fi

mkdir -p "$INSTALL_DIR"
echo "📂 aperfetch '$INSTALL_DIR' dizinine kurulacak."

echo "⚙️  Çalıştırılabilir dosya kopyalanıyor..."
cp "target/release/aperfetch" "$INSTALL_DIR/aperfetch"

echo "chmod +x "$INSTALL_DIR/aperfetch" izni veriliyor."
chmod +x "$INSTALL_DIR/aperfetch"

echo ""
echo "🎉 aperfetch başarıyla kuruldu!"
echo ""
echo "Şimdi 'aperfetch' komutunu terminalinizde çalıştırabilirsiniz."

if [ -z "$TERMUX_VERSION" ]; then
    echo ""
    echo "Eğer 'aperfetch: command not found' hatası alırsanız,"
    echo " '$INSTALL_DIR' dizinini PATH ortam değişkeninize eklemeniz gerekebilir."
    echo "Bunun için aşağıdaki satırı kabuk yapılandırma dosyanıza (~/.bashrc, ~/.zshrc vb.) ekleyin:"
    echo ""
    echo "  export PATH=\"$INSTALL_DIR:\$PATH\""
    echo ""
    echo "Değişikliklerin etkili olması için terminali yeniden başlatın veya 'source ~/.bashrc' (kullandığınız dosya) komutunu çalıştırın."
fi

echo ""
echo "Mutlu kullanımlar dileriz!"

