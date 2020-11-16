package main

import (
	"fmt"
	"github.com/limechain/hedera-state-proof-verifier-go/stateproof"
)

func main() {
	bytes := []byte(
		`{
    "record_file": "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",
    "signature_files": {
        "0.0.4": "BB7i2u97gqhFhXziRcuVB1kbjaNmp7zH8YkfJEK18vJ1VOPdaqQJAPBIF4zz9P+NkwMAAAGANg7fnHr38/zs1auH7+nAe/TI0E0UZmGZiRgJqGIQRj1ndnqy4a2u01C9Ycu9R2SsFI89d5TzT9F2MHr74qecLrwVm36ZNwaX29AQcgM/JGDVtzXZ82rwszBDEHsYPAl/xzZ6FoTjEM9nlaiKW3KpMZB8XrhxQvJE1ekwvGeWVAMJF0/KT3G7zUjC0GZBmoxlI5/z3DvfsAP1EVeSW/zoAUJQUqitNKteRY0lGZJgfUJBHqboJYwt03CulDedzGaptiZ0LhL5IOj9eeTl/0EekLlvDp4Q81Z6hSMTd1zxZMnVNb7/vj8Z+u7wU0fLofyge4mdl8fn2E+CED1N/bZRNZlZvAB5x3SfvWA5DlKG4CvyEWt2bcDywrxnKLskGsSf5+hzL0bC/NOs7jog8leOGw515XT91Rh2Te5NnnKw3EP8oYnD5Wlgh9Gme+t3h1yDn4Nv7euJP6L4oBnY5G80EmkhQyZFRwpN7j881GBcUedh37U9ZcaBtjt8AmBpwTn4",
        "0.0.3": "BB7i2u97gqhFhXziRcuVB1kbjaNmp7zH8YkfJEK18vJ1VOPdaqQJAPBIF4zz9P+NkwMAAAGAEfSanSB4KaG+YQG1chLG1IwgfGbLr2aFiMvSMZInW7nF0/cU/tKc8ygp8WgiQTV8G4wPhzwPtmtnb4wKCPks0SEocdLsnIh3e9ugh1MDjgaKo+T3YfBLV72mn7iDlkg+MABJ4FttptG+IH4zQrDfXEbyTf/VqIC4CEnxXzuHE9vO6YILfrzEAZBatxvz9K2Jqw2rWfSAWkicpRqkIqEDL+v9In3mw8oU0DGC1imty1As9P0DLsMbRE3jueaqYdR1kPp5sDlrrAC7BpcaYP0LjCZB25bLXsAtsh3NAYbtz5y3YxdV2wDpakN2pyT8W6c+QEPnGwRQfSGdW9Meobd5SMGtltJ8u/IHTqurWJllj3nWJd7juo3ms7S2YwbIJNDqGUuLJxWHxlT5RgJn3IGRc4IT5ROtRzyudtUyWYUBDRGKA9QsQ28WCq/UNqIzc3t4pG0dJ9/ucwKALHcBgAa0fK+ta3kjQhFAsFzD3k7apC/OsEzoFkunsGJY5uIe9iLb",
        "0.0.6": "BB7i2u97gqhFhXziRcuVB1kbjaNmp7zH8YkfJEK18vJ1VOPdaqQJAPBIF4zz9P+NkwMAAAGAE1zHNbw2qZjjQZCPdLGVrlPqnZPryiTksTwoU8FQs6wiJp0DnX96yoaWpNlk25fb+nNrYWMi164FWSNSdFBxo4megiWzKpAcbVZJefiCVaqssDfL3o4Xti/RMKvjdU1XFyqv3OyiHym4XldFP0tPhZNNpyhZxjwuLVvqtBN8BXXTZYfcREcZ9AzZtc2AYOuXKel7zN2hMKJ5+Nz09hOkaOiMsKCNjPEQLMCslvZaJS/tyO3I2WNG0PaUmvkf3VHxkqzTmCqhV/bReV5YZ5iIcsEok5/Vdu/9iLrW/aRTuqksrrP2N24yKAsAlkHFBW/GFzx/LKDXW9ifyQzq+n1U5zm4zJc9VaPprI+aN9hgYiqU+dHFiUvGLA3F5EamAf6nqJGHXX0Req3u9b0tBg7SA55YHpTEitKLC/sklsfzGPbrigBzRuhZvqVBGiHAHYJ7uNtNH/W2muz68tdDvoGG4YqYqDJZQqhfE/L4wZtx/i65r3LBgNsQtmz+9oAdv838",
        "0.0.5": "BB7i2u97gqhFhXziRcuVB1kbjaNmp7zH8YkfJEK18vJ1VOPdaqQJAPBIF4zz9P+NkwMAAAGAFbKZLQkMBzHRbs9AYgYNgflGt+I4AECU28f+VxebziNAlgv86YNvRvMsZoVJjVndxZcoWPRzuszPdHohH+tyw6I4mL58Qk1dJ0fE2F2TAt88V+PpN+AUAR6HopwQfbKhZzdnxDh4EOCqullYsfVhIaBk48WEuuoGQvUhTxFbx8gbnCaUXccVa3AwlkVYqFheGNymTwIDLFHq+LCIHIdnVba95FFhAvqO1b+1ARSd2ZjnbvOJ5QaEoffQ+Liza6SgQ5/UsIogC+70z3sfnxZU0SW0dRbnqvkYpBxwKmTkiYP4CGPIr1CyEAyh1hRJdZxXIhti74fi4i+jiDvpsHdcS4qtREn9V/pb5/EayckBhiO4sslHdKaLCyIvot1W4zcHcQT8RP0dGetSTJ4NNiIeSbEo6ouLxThcMxGgEGrAYU8tvGQoIqnMXl5W5GpmxNgHnOHrM87nFv68FdYRoBFGs061KpdgBZo4/NFfhC60rYmaRCX7PmNvS8CQtmZXQyLG"
    },
    "address_books": [
        "CtYGGgUwLjAuNSLMBjMwODIwMWEyMzAwZDA2MDkyYTg2NDg4NmY3MGQwMTAxMDEwNTAwMDM4MjAxOGYwMDMwODIwMThhMDI4MjAxODEwMDliYTQ1N2I3MzMwNWYwNGE5MWNjNDZiMWI5NjVjNGU4NDE3NTFhYmM4YjE0MTVhMGJhZGZkMWYzMmMyNDgyMzg2YTIyNzI1ZWI3ZWM3NGRlYTIxZTUwNjE3ZDY0OGVhNWFjMzkzNzQxYWIwMWI4ZWZiMzIxMjM5YjhkNGZkYjFkZmJlYjllM2YzOWFhNDY1ODBkZDA0NWQxOGNhNDRkMDAyYzM3ZGRiNTI3Y2NlNGRkYzMyYmZjNzM0MTk2NzFmNGNhNDQ2NGEzZjJhODRmYzg1YzcxYWNmMGU1YTg5NjI2ZGY2OWE4MTQ3NGVkMTY1MjlmODAxYThhZmE5N2U0MzVjNGUwNGE5NjRhMzU3NTI3Mjg4ODQzZTU4ZjBhMDVjZjUxNTNlZTQ1MDdiMmM2OGIzZDdmYjU0YWU2YTk1YTk1OWM4N2ExMmY2MzBlOTVjN2IxYjNjMzY5NWU4NTg2NjI0MTc5MjZkNzZjMTY5ODNmYWY2MTIyNTAzODc0NTkwN2U5Y2YxM2Q2N2MyYWNkNTAzY2E0NTFjODU5MzNhYzQxMThhY2MyNzk4MDFjYjk2ODM0OTkwMzE0NWNlZDI3NjI5ZGQwODkxNjMxNzA5MzU4N2E3N2MyMjA1Y2ZhNTI1NDNiNTNjM2I2ZWExNWI4NGUzZDJjMzBjMWVkNzUyYTQ2MzNjMzZiMjViOTg5M2VhMDJhZDU2MmViOWI3ODY4YjNiNGY0N2Y0YTI1ZTM1NjA2NDk2MmFjN2IyNWU1ODI5NDRmMDBkMzA3OThhMjYyZjkyMTRkOGM1ZTc0ZDBhODM3NmNjMmQ2YmE2NGUxOGY1ZTRhNDBhZmFjNjI1MDYyZDJjYTIzY2QyODAwNzA4MzIxZDM4MzQzMTRmMGU1ODQ0ODU5MjMyNjczYTMyZTcwYWUwZDcxMWUzMTA1ODFiY2RiMTRlODcxMzQ2OTRjNmUwOTMwZjQ2YjM3Yjk2ZDQ5YTY0NTczOTQ3MzMxZTdlNTA3ZDllNTZkZTVlNjE0NmYyZjAyMDMwMTAwMDEK1gYaBTAuMC42IswGMzA4MjAxYTIzMDBkMDYwOTJhODY0ODg2ZjcwZDAxMDEwMTA1MDAwMzgyMDE4ZjAwMzA4MjAxOGEwMjgyMDE4MTAwYzQyY2NhYzVmYmM2OTFmYmJlYmRhODdmZmQxZTc1YmRjZDg5MjI0OTRjZjQ0ZmRiY2NlZTQ5Nzg4NTIxYzM3OGJmNzdkYjA5MzRlYzBkMjE4M2Q3YzUxZGI2NmY4NjRjMTFhYjdkZTFhYzNjNGNmZGMxZjA5M2EyZDZmMzdlMmIzNGNiZTRjODEzMWY5NjgzYWQ0Mjg3OGM4M2QzNTU0YzY0NWFhMTY3YmNmYjA2NGE4M2RjNDVjNWIxMTU4NDk5ZjlkOTI1ODdmZmY3YWJjZDVmMjIxY2Q4MTUwNTQ4NDEzMDAwZmE2ZTU2NTkwODliMWRmZDY1NzY2ZWE3OGVhZWRmY2E2YjQ1NDU1ZmQ4YWI1OTg0ZGJlMzVlNTc5NWQyYzYzNWVhNzk3NGQ0M2U4ZWFlNGZlYmZmZTQ5MmU3MDdiNDhiMWIwZmM2NDgxYWU5ZTA5ZDM5MTMzMDA5YjdkMjY0MDJlNmU1MmU1ZTkxYjJiMzgwZDg4ZjBiZTdmYjRiMzAzZTcwMjE5Nzg1MDU3YWE5NGNlOTI0YzQ5MjZlOTE2NTY5Mjg2ZTg2YjNiYTY1MWNhMmEwYTYzZGY0ZjY5MDdmZWZlMzQ4M2Q5M2I0Y2UxZDRkMDNjNzE0MjExMTM3NWIyYzJjNTFkNGViODM5ZTM3YWY1MzBiMmNiZDZmNTBkNGNiMzZlMjc5MzcxNzBkOWNkZGFjMGFjZTJjYzI0YjgwNGIwYTI3MzUxY2Y4MzBiNzY1MjVlMjZkZmI5ZGJmNDlhMDU2NjI0YTc2ODYyNDk0ZTcyNjNkMGQ3MGNlYmFlOTUyOTQzZTU1ODQyZjVjYWQxM2ZjZjYwYTJlNmRjZjdhMWQ1MzNmM2E1YmI1NGVjMjE5MThjNzZlNTI1YmEyOTE0NjY3NTgzMWUxN2UzNmM2MWZlODU0OTg4MjhkMDliNzYyMDE1NDEyYjJlNTI3ODQ5YmFlYzFjZmZjNzdkZTRjMjk0YzU1MDgxMWU1OThmZjI0ZGExNWEzNDU2OWRkMDIwMzAxMDAwMQrWBhoFMC4wLjMizAYzMDgyMDFhMjMwMGQwNjA5MmE4NjQ4ODZmNzBkMDEwMTAxMDUwMDAzODIwMThmMDAzMDgyMDE4YTAyODIwMTgxMDA5ZjFmOGExMjFjMmZkNmM3NmZkNTA4ZDNlNDI5ZjBjNjRiY2I0NGM4MmE3MDU3MzU1MmFhZGNhZDA3MTU2OWU3MjE5NThmNWE1ZDA5Zjk1ODdmZmFmY2ZiZTUzNDFhMmYwMTE0YWNhZTM0NmVmM2M5MDIxM2QzNDM2ZWJiMjdmNDM1MGM5OTBjNWM4YzNmOGUxZTM2NzA3YmMwOGQ0MjU2MDgyM2UzZjI0ZTA5YTAzYWQwOTU1YTUwOTgwMTk2MjlkZDA0YjI3YjI1MWRjZTA1NWYzZGRjYjBhNDFkNjZmMDk0MWIwYjg3Y2RmZTM0OThkNDYwMzhhYjVkZjA2ZjYyYTVhZGUwODU5ODU3M2E4OGM4ZjU4NjBkYzE0OTJhNmUxODY0ODVhOWIxMzI1MGU2ZDE3YjgwY2QzOWM1YzgxOTEwOWU3M2NhNzMyZGIyM2VmOGJhYTc3NmVjODVjZTAwOTFiZWNiMmVkZWZiYWE1ZWQzZTVkYmZiZDFmODg1YTRmYTg4MWFmM2YxNDRhOGE1NjU4NTM1MzNkODkzOTM1OTIwODZiMmQxZDM2MmU0NWJmZTFmYjQ1NjgzYWJhNmM2NDA5NzlhZDZiNDY4NzcxODQ3MjZjNmViZDU4YjJlYWU4NWM3Y2ZlM2ZiYWJlZjVmNmNjZWQ4NTAwMzRiMzg0NzIwNmMyZDY3OGMzNjE4NzYwMjZiOGQzNTFlMDAyYWY1ZTBmZmU2ZjViMWYyOTVmZGMyZjQ2OWNhYTJkMjM4MWVhMGI0OGNhOTg3Y2MyYzhlNjM1ZThiMTljZTVlMTcyYTkzNzYxYThkNDkwYTlhNDUxOGQ3MjU1ODgwYTE0ZDc3YjdiYTc3NDg5MmI5MmE0MGJiODEzNjJlMzRmYzZkNTE3OGQ5YjMwMTEyOTM0MjA1Y2I3N2ZiOWEyODI0MjczOTQ1NjRhODU1NGVhNDcyODZhNDdmODYyMzllNzVjOTQ3ODljZTk4Yzk5ODQ0NzgyNDYyOTQ0ZjYxMzE2N2Q3YjUwMjAzMDEwMDAxCtYGGgUwLjAuNCLMBjMwODIwMWEyMzAwZDA2MDkyYTg2NDg4NmY3MGQwMTAxMDEwNTAwMDM4MjAxOGYwMDMwODIwMThhMDI4MjAxODEwMGM1NTdhZjU3OWZhODM1MDFiZTg5OWIyODkwNzc2NWJmZGZjZDUyYWI0MzJiMDE5NWExZjFlY2Q4NmZjMDBhYjZjNTUwOWIwZmRkOTdlZGQzY2I1Y2VhNTZhMjk1ZjMxMmFiYjU1MDgzMWRiZjk2M2Y0NTAxMThiNGZjYzZlMjJjZjQ2NzYyMDBjZTljYzhlZGZiYmY1NThkYzY5ZjAyNDI2NGFkN2QzZGFiMjNiZWQyMTMzYzI3NGU2OTM0NDg5MTU1ZGIxMDg3ZjkwMzcwOTA1YzY0MTg1YTYyMTFkYzc0MmZiOWE2OTA5ZDgyMTg2OTQ3YjI3NzQ2M2RmYjNmZjBhY2Q0N2VmZjEyZWFkMWY2OTcyZWYyYzEyMDM3OTNjNDVlNzc1NzViZTRmYTExMGM3ZTQwZmE4ZGI5YzYxODdkMTEzZjQ3MDQwMTQxNzkwNzFhYmY1OWJlN2QyYjBkZTgyZGU0MjE1ZGMyNTUwNmIxYzljMjZlNDkxNzQwMWM5OTc1MDZlMzc3ZTZiZjAzYjY4ODcyN2U3OTQwZmFkNjljNWUwZGEzY2Q1Y2JkMmJlNzc3MzUwYWVhMmQwZDQ3ZTk3YTQ0OGM4NGJlNmNlMTM0ZDY0YmVlMDk4NWMyOTE2MmY0YzFlNTY3Y2NhOTNkMDZhM2MxYmU4YWJjZTM1YjU1N2ZiNzdmNGZlNjcxYTY2ZGVjNzkwNzU2ZDBlODgxODE2NWYyYmFjYWE4OTFhYWU3YWM3NDM3ZmM3MTc1YjZlYjZkZWI3NDcyMzc4NzUxYmI2YmY5YjBlMTQ4M2Y5NjY4ZTlmZGJkNTYwNGMzOWIxNGQ5ZTJiZWRlZWM4NDZhOTgwZDcwNGQxNzFlN2JhNGI3ZmNkMWEzMGQ5NDVjYTEyZjQ3YTMyNWQ5Mzk4YWExOGY5NzA2NjA1NGQ0ZDE1ZmM4OTk0ZTJkZWJlNzNlOTI3MWQ1NDg2ODNmNjFlYTQ0ZmIyNTA3MWUzNTE4YTc4ZWQzZWIzN2U3MWEwNjkxZjI2NzAyMDMwMTAwMDE="
    ]
}`)

	res, err := stateproof.VerifyStateProof("0.0.1893-1605177623-307000000", bytes)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}
