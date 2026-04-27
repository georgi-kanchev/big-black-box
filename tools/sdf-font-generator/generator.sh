#!/bin/bash

# --- CONFIGURATION ---
RESULT_DIR="./results"
PACKAGE="msdf-bmfont-xml"
CHARSET_FILE="$(pwd)/.charset.txt"

# Define the character sets
CORE=" .,;:!?¡¿\"'()[]{}<>-/\\@#$%^&*_+=|~\`0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
LATIN="ÀÁÂÃÄÅÆÇÈÉÊËÌÍÎÏÑÒÓÔÕÖØÙÚÛÜÝßŒŠŽŁŃŚŹŻĆČĐŐŰàáâãäåæçèéêëìíîïñòóôõöøùúûüýÿœšžłńśźżćčđőűẞ"
CYRILLIC="АБВГДЕЁЖЗИЙКЛМНОПРСТУΦΧЦЧШЩЪЫЬЭЮЯабвгдеёжзийклмнопрстуфхцчшщъыьэюяҐЄІЇґєії"
ALL_CHARS="${CORE}${LATIN}${CYRILLIC}"

echo "--- MSDF Font Generator (XML + No Padding) ---"

# 1. Global Dependency Check
if ! npm list -g $PACKAGE --depth=0 &> /dev/null; then
    echo "📦 Package '$PACKAGE' missing. Installing..."
    sudo npm install -g $PACKAGE || { echo "❌ Install failed"; exit 1; }
fi

# 2. Input Handling
read -p "Drag & Drop TTF file: " FONT_PATH_RAW
FONT_PATH=$(echo "$FONT_PATH_RAW" | sed "s/['\"]//g" | xargs)

if [ ! -f "$FONT_PATH" ]; then
    echo "❌ Error: Cannot find file at '$FONT_PATH'"
else
    read -p "Enter font size [48]: " FONT_SIZE
    FONT_SIZE=${FONT_SIZE:-48}
    FONT_NAME=$(basename "$FONT_PATH" | cut -d. -f1)
    
    # 3. Create Folder and Save Charset
    mkdir -p "$RESULT_DIR"
    sudo chmod 777 "$RESULT_DIR"
    echo -n "$ALL_CHARS" > "$CHARSET_FILE"

    echo "⏳ Generating MSDF Atlas..."

    # 4. Execution
    # -p 0: No padding between glyphs
    # -b 0: No border at texture edges
    npx -g $PACKAGE \
	  --pot \
	  --square \
	  -t psdf \
	  -m 512,512 \
      -f xml \
      -s "$FONT_SIZE" \
      -p 0 \
      -b 1 \
      -i "$CHARSET_FILE" \
      -o "$RESULT_DIR/$FONT_NAME" \
      "$FONT_PATH"

    # 5. Rename .fnt to .xml
    if [ -f "$RESULT_DIR/$FONT_NAME.fnt" ]; then
        mv "$RESULT_DIR/$FONT_NAME.fnt" "$RESULT_DIR/$FONT_NAME.xml"
    fi

    # 6. Cleanup and Permissions
    [ -f "$CHARSET_FILE" ] && rm "$CHARSET_FILE"
    
    sudo chown -R $USER:$USER "$RESULT_DIR"
    sudo chmod -R 666 "$RESULT_DIR"/* 2>/dev/null

    if [ -f "$RESULT_DIR/$FONT_NAME.xml" ]; then
        echo "--------------------------------------"
        echo "✅ Success!"
        echo "Format: XML (No Padding)"
        echo "Files: $FONT_NAME.xml, $FONT_NAME.png"
    else
        echo "--------------------------------------"
        echo "❌ Error: Generation failed."
    fi
fi

echo ""
read -n 1 -s -r -p "Press any key to close..."
echo ""