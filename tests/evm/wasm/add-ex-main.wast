
(module
    (import "debug" "printMemHex" (func $printMemHex (param i32 i32)))
    (import "debug" "printStorageHex" (func $printStorageHex (param i32 )))
    (import "ethereum" "getCallDataSize" (func  $getCallDataSize (result i32)))
    (import "ethereum" "storageStore" (func $storageStore (param i32 i32)))
    (import "ethereum" "getBlockDifficulty" (func $getBlockDifficulty(param i32)))
    (import "ethereum" "callDataCopy" (func $callDataCopy(param i32 i32 i32)))
    (import "ethereum" "getCaller" (func $getCaller(param i32)))
    (import "ethereum" "getExternalBalance" (func $getExternalBalance(param i32 i32)))
    (memory 1)
    (export "memory" (memory 0))
    (export "main" (func $main))
    (func $main
        ;;-------------------------
        ;; 变量赋值
        ;; dataSize = getCallDataSize()
        ;;-------------------------
        (local $dataSize i32)
        (call $getCallDataSize)
        set_local 0

        ;;-------------------------
        ;; 变量值，写入内存
        ;; memcpy(ptrMem+0 , &dataSize, 4 )
        ;;-------------------------
        (i32.store (i32.const 0) (get_local $dataSize))
        (call $printMemHex (i32.const 0) (i32.const 100))

        (if (get_local $dataSize)
            ;;-------------------------
            ;; dataSize 不为0，说明有入参
            ;;-------------------------
            (then
                ;;-------------------------
                ;; 第一个入参(opType)的值，写入内存
                ;; memcpy(ptrMem+4 , &data+0, 4 )
                ;;-------------------------
                (call $callDataCopy(i32.const 4)(i32.const 0)(i32.const 4) )
                (call $printMemHex (i32.const 0) (i32.const 100))

                ;;-------------------------
                ;; 第一个入参(opType)， 如果是 0x01000001， 则进行 mint 操作
                ;;-------------------------
                (if (i32.eq (i32.load (i32.const 4)) (i32.const 0x01000001 ))
                    (then
                        ;;-------------------------
                        ;; 铸币地址，写入内存
                        ;; memcpy(ptrMem+8 , &callerAddress, 20 )
                        ;;-------------------------
                        (call $getCaller(i32.const 8))
                        (call $printMemHex (i32.const 0) (i32.const 100))
                    )
                )
            )
        )


;;      (call $printMemHex (i32.const 11) (i32.const 32))
;;      (i32.store (i32.const 11) (i32.const 101))
;;      (call $printMemHex (i32.const 11) (i32.const 32))

;;      (i32.store (i32.const 0) (call $getCallDataSize))
;;      (call $printMemHex (i32.const 0) (i32.const 32))

;;      (i32.store (i32.const 0) (i32.const 0x01000001) )
;;      (call $printMemHex (i32.const 0) (i32.const 32))

;;      (call $callDataCopy(i32.const 0)(i32.const 0)(i32.const 3) )
;;      (call $printMemHex (i32.const 0) (i32.const 32))

;;      (call $getCaller(i32.const 4))
;;      (call $printMemHex (i32.const 0) (i32.const 32))

;;      (call $storageStore (i32.const 100) (i32.const 0))
;;      (call $storageStore (i32.const 11) (i32.const 12))
;;      (call $printStorageHex (i32.const 11))
;;      (call $getBlockDifficulty(i32.const 11))
;;      (call $printStorageHex (i32.const 11))
;;      (call $storageStore (i32.const 0) (i32.const 0))
;;      (call $printStorageHex (i32.const 100))
;;      (call $storageStore (i32.const 100) (i32.const 0))
;;      (call $printStorageHex (i32.const 100))

    )
)
