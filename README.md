# scalable lightweight accounting go api

# Windows

    - mattn/sqlite3 requires gcc compiler for arch, if on windows and x4 you need TDM_GCC_x64

# TODO  

- [.] Misc
  - [x] centralized logging w debug mode
  - [ ] log to file
  - [ ] Actual migrations sytem/tool
- [ .] Accounting
  - [.] Transaction CRUD Impl
    - [x] POST
    - [x] GET byId
    - [x] GET ALL
    - [x] UPDATE
    - [ ] DELETE (virtual delete)
    - [ ] Transaction schema fix
  - [ ] Transaction CRUD Integration Tests (httptester)
  - [ ] balance(start_date, end_date, filters) Impl
  - [ ] balance func testing
  - [ ] balance func usages Interface
- [ ] Auth
  - [ ] User CRUD
  - [ ] Role CRUD
  - [ ] Authentication with OAuth
  - [ ] (Store)Role based Authorization
  - [ ] Auth middlewares
- [ ] Ecommerce
  - [ ] Store CRUD
  - [ ] Product CRUD
  - [ ] Cart CRUD
  - [ ] Order CRUD
  - [ ] Store POV accesses
        - [ ] role based Write(CreatUpdate) permisions on products by store
        - [ ] created (published/draft/private) products by store
        - [ ] created (pending/fullfilled) orders by store
    - [ ] Buyer POV accesses
      - [ ] Available stores (store ranking/reccmendation system)
      - [ ] Product & Store FullText Search
      - [ ] saved carts by user
      - [ ] created (pending/fullfilled) orders by user
- [ ] Payment
  - [ ]v0 "CONTACT SELELR AFTER CREATING ORDER"
  - [ ] Payment Gateway Provider
  - [ ] compliance of api
