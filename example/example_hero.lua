local t = {
    [1001] = { 
            id=1001, 
            name="陈咬金", 
            atk=100, 
            def=120, 
            mp=12
    },
    [1002] = { 
            id=1002, 
            name="孙悟空", 
            atk=300, 
            def=200, 
            mp=300
    }
}

require'metadata'.new((...), t)
