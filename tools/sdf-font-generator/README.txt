make sure to fit the font size into a single 512x512 texture!
- having it too small will shrink the texture to 256x256 or lower
- having it too big will generate a second page texture
only a single 512x512 texture works with the engine!

this is because the atlas will carry its glyph data on the very first row of 512 pixels
each glyph info fits into 2 pixels and the charset is a forced 256 symbols x 2px = 512px