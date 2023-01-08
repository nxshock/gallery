# Gallery

Simple gallery.

## Usage

1. Create config file and set `WorkingDirectory`
2. Start program

   ```sh
   gallery <path to config>
   ```

## Config structure

```toml
WorkingDirectory = "/home/user/pictures" # Working directory
Crf              = 40                    # Preview CRF
                                         # from 1 - maximum quality
                                         # to  63 - minimum preview cache size
ProcessCount     = 4                     # Preview generator threads
```
