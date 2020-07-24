package main

import (
	"cpabse"
	"encoding/gob"
	"fmt"
	"os"
)

func main() {
	/*
		pm := new(cpabse.CpabePm)
		msk := new(cpabse.CpabeMsk)
		cpabse.Setup(pm, msk)
		bpm := cpabse.PmToBpm(pm)
		filename := "pm.txt"
		f, _ := os.Create(filename)
		enc := gob.NewEncoder(f)
		_ = enc.Encode(bpm)

		bmsk := cpabse.MskToBmsk(msk)
		filename1 := "msk.txt"
		f1, _ := os.Create(filename1)
		enc1 := gob.NewEncoder(f1)
		_ = enc1.Encode(bmsk)
	*/

	bpm := new(cpabse.BytePm)
	f, _ := os.Open("pm.txt")
	dec := gob.NewDecoder(f)
	_ = dec.Decode(&bpm)
	pm := new(cpabse.CpabePm)
	cpabse.Psetup(pm)
	cpabse.BpmToPm(pm, bpm)

	bmsk := new(cpabse.ByteMsk)
	f1, _ := os.Open("msk.txt")
	dec1 := gob.NewDecoder(f1)
	_ = dec1.Decode(&bmsk)
	msk := new(cpabse.CpabeMsk)
	cpabse.BmskToMsk(msk, bmsk, pm)

	attrs := "foo fim baf"
	prv := cpabse.CP_Keygen(pm, msk, attrs)
	policy := "foo bar fim 2of3 baf 1of2"
	keyword := "love"

	c, cph := cpabse.CP_Enc(pm, policy, msk, keyword)
	fmt.Println("------------------------------------------")
	fmt.Println(c)
	fmt.Println("------------------------------------------")
	fmt.Println(cph)
	s := "54 255 129 3 1 1 7 66 121 116 101 67 112 104 1 255 130 0 1 5 1 1 87 1 10 0 1 2 87 48 1 10 0 1 2 87 49 1 10 0 1 1 83 1 10 0 1 1 80 1 255 132 0 0 0 124 255 131 3 1 1 10 66 121 116 101 80 111 108 105 99 121 1 255 132 0 1 10 1 1 75 1 4 0 1 4 65 116 116 114 1 12 0 1 1 67 1 10 0 1 2 67 112 1 10 0 1 8 67 104 105 108 100 114 101 110 1 255 134 0 1 1 81 1 255 136 0 1 11 83 97 116 105 115 102 105 97 98 108 101 1 2 0 1 10 77 105 110 95 108 101 97 118 101 115 1 4 0 1 5 65 116 116 114 105 1 4 0 1 4 83 97 116 108 1 255 140 0 0 0 35 255 133 2 1 1 20 91 93 42 99 112 97 98 115 101 46 66 121 116 101 80 111 108 105 99 121 1 255 134 0 1 255 132 0 0 46 255 135 3 1 1 14 66 121 116 101 80 111 108 121 110 111 109 105 97 108 1 255 136 0 1 2 1 3 68 101 103 1 4 0 1 4 67 111 101 102 1 255 138 0 0 0 23 255 137 2 1 1 9 91 93 91 93 117 105 110 116 56 1 255 138 0 1 10 0 0 19 255 139 2 1 1 5 91 93 105 110 116 1 255 140 0 1 4 0 0 254 6 124 255 130 1 255 128 91 211 236 241 20 23 96 212 104 23 226 56 91 196 217 128 5 58 17 71 81 141 198 230 100 42 38 120 118 197 190 48 52 42 176 165 42 72 48 99 172 252 111 181 13 153 145 54 113 187 228 60 88 70 67 127 40 96 32 196 246 108 113 8 149 12 64 3 115 111 239 70 211 62 6 59 223 206 93 104 28 206 164 219 234 192 108 84 35 241 255 144 159 96 223 85 175 66 145 103 216 123 24 52 174 132 22 70 106 204 29 93 248 208 36 231 83 84 147 3 87 5 142 163 133 215 71 233 1 255 128 26 69 32 87 88 165 19 164 19 146 168 127 89 147 135 181 211 132 231 41 33 71 205 194 245 85 90 161 74 228 32 115 225 120 218 18 236 14 13 168 93 216 190 198 13 138 24 40 155 79 28 206 221 40 72 38 32 140 204 140 201 25 131 83 133 173 133 11 59 68 188 149 61 55 80 225 213 222 141 45 54 69 6 250 153 236 106 201 22 53 247 9 36 97 96 128 36 20 118 192 74 192 144 80 68 20 201 26 32 44 52 52 149 144 139 17 27 206 216 234 108 123 113 75 72 86 47 173 1 255 128 83 148 34 0 16 121 166 195 213 186 37 108 233 82 104 53 241 42 42 43 129 95 63 102 19 193 7 53 163 6 173 47 19 72 91 133 6 181 127 59 80 204 169 153 2 202 222 253 136 101 11 231 146 143 212 138 146 91 115 230 83 95 89 132 53 251 129 76 109 218 179 44 218 157 68 200 43 250 239 201 84 110 182 72 105 186 9 143 112 49 222 98 23 170 223 224 190 113 129 85 73 31 246 198 117 65 178 85 109 218 52 124 254 88 142 145 145 140 90 215 138 102 19 102 114 112 80 195 2 1 2 4 2 1 4 4 3 1 2 1 3 102 111 111 1 255 128 153 164 253 154 250 185 131 243 101 204 208 213 173 184 88 194 83 21 33 140 125 79 40 226 166 127 123 117 27 167 156 96 43 147 103 95 142 46 127 251 179 54 93 15 117 127 12 93 83 0 192 153 100 154 38 193 124 247 243 237 81 168 69 250 21 19 171 159 199 104 198 100 186 43 126 133 92 111 46 132 41 174 243 20 110 88 100 196 44 205 116 247 36 252 73 160 117 156 85 63 100 185 7 7 223 159 81 50 82 245 88 77 86 63 121 34 249 55 67 89 173 104 101 139 220 68 49 14 1 255 128 155 116 68 63 185 112 210 104 61 43 209 213 147 244 58 208 89 108 26 101 33 103 241 32 130 173 110 55 166 67 17 195 154 163 182 63 11 111 246 204 249 224 83 8 238 2 39 150 24 7 19 143 187 202 178 162 137 214 208 58 122 148 248 154 0 167 152 105 46 86 18 114 255 150 190 175 169 58 197 150 10 38 106 8 97 211 104 236 148 11 219 164 161 249 174 139 234 65 195 250 134 9 222 48 117 173 219 168 179 102 55 60 13 178 251 23 227 198 161 106 248 48 133 31 123 180 73 254 2 2 1 20 32 103 54 74 129 223 5 158 225 179 214 170 41 79 9 213 150 235 42 231 0 0 1 2 1 3 98 97 114 1 255 128 141 161 244 168 48 14 206 200 116 193 168 99 8 214 241 179 40 12 172 165 25 6 246 107 55 207 129 219 121 120 109 10 149 206 74 234 107 193 134 66 235 8 91 39 185 97 199 172 203 16 79 36 78 64 246 66 205 12 211 37 86 103 199 158 58 167 237 142 186 9 87 75 255 251 199 117 29 243 248 134 241 235 205 174 169 131 177 39 108 57 9 231 193 0 179 223 20 253 137 23 255 26 30 240 195 148 208 171 215 141 25 15 245 180 112 177 59 245 6 27 25 31 187 30 108 225 42 5 1 255 128 147 237 15 158 59 85 221 47 64 243 89 146 23 242 53 164 188 93 34 52 53 197 30 211 106 236 39 227 21 160 55 32 107 63 247 231 157 245 173 99 77 190 60 46 248 122 77 60 149 85 7 41 68 210 204 57 90 234 214 101 123 61 31 240 148 40 98 112 139 61 40 143 130 103 195 164 223 36 125 41 88 195 219 250 164 39 176 247 233 140 154 63 109 138 66 142 19 209 116 37 181 252 50 161 73 95 30 155 160 254 2 246 37 240 225 253 222 177 135 162 47 55 73 67 254 149 1 84 2 2 1 20 40 105 102 8 162 28 104 139 71 248 72 173 181 246 24 64 14 63 158 187 0 0 1 2 1 3 102 105 109 1 255 128 17 121 60 53 155 107 196 122 18 30 112 24 15 189 94 142 15 147 232 172 111 188 193 138 196 113 19 240 245 30 119 59 78 27 209 50 199 40 182 24 216 235 186 146 155 118 48 88 222 243 204 21 78 171 127 43 177 19 18 115 218 202 224 205 157 15 247 122 17 70 7 81 35 226 18 82 230 184 27 126 71 15 226 38 211 70 80 175 133 112 243 162 223 185 11 64 15 69 61 227 77 95 222 118 71 200 1 242 90 82 108 44 215 0 92 141 62 73 118 81 137 113 128 63 163 237 217 215 1 255 128 18 18 140 65 13 209 231 136 0 255 52 248 40 55 93 27 179 84 248 38 215 248 207 188 43 115 99 124 212 37 153 201 156 101 249 140 210 23 31 93 4 79 238 201 71 194 123 164 204 155 23 27 69 113 81 137 136 5 218 208 211 222 209 206 14 236 185 192 243 124 118 53 125 67 64 90 126 134 91 49 198 155 15 210 231 120 144 240 158 205 22 35 21 142 141 121 212 108 40 178 41 82 235 198 251 24 214 226 102 225 25 53 218 159 229 5 97 45 189 123 42 101 40 18 174 87 128 144 2 2 1 20 48 107 149 198 194 89 203 119 174 60 186 177 66 157 38 170 133 148 18 143 0 0 1 1 2 1 2 20 24 101 6 140 97 161 162 178 123 111 100 166 156 167 251 107 31 150 183 19 20 8 2 47 190 32 61 98 236 102 68 114 3 140 167 14 106 119 84 115 212 0 0 1 2 1 3 98 97 102 1 255 128 124 58 116 17 38 50 238 11 152 144 216 200 45 215 136 174 159 26 149 6 118 247 56 235 43 53 180 159 95 221 212 95 158 212 75 202 95 202 231 199 53 154 67 203 201 41 212 84 232 183 151 129 85 73 237 15 239 34 144 106 171 47 104 33 28 152 239 7 30 23 90 48 247 4 119 246 232 142 97 189 131 146 56 53 3 254 42 55 228 67 220 44 187 155 99 180 115 129 160 101 160 183 80 50 196 184 38 143 187 100 80 193 203 68 247 241 20 131 92 235 138 178 104 76 122 66 6 195 1 255 128 25 55 216 77 123 111 72 62 37 162 125 246 11 7 178 67 113 115 43 190 50 45 179 37 33 14 12 249 41 29 225 214 119 194 229 204 165 118 47 147 185 32 2 5 237 19 208 28 255 147 123 226 124 201 216 95 220 21 84 138 181 201 37 162 158 251 103 157 221 75 44 38 219 65 12 25 55 43 242 165 204 7 223 224 224 165 156 216 34 17 142 162 211 75 43 167 232 25 64 86 48 108 201 227 195 168 3 220 242 159 59 100 77 217 230 183 140 169 234 144 46 63 87 204 96 40 152 161 2 2 1 20 24 101 6 140 97 161 162 178 123 111 100 166 156 167 251 107 31 150 183 19 0 0 1 2 1 20 24 101 6 140 97 161 162 178 123 111 100 166 156 167 251 107 31 150 183 19 0 0 0"

	cph1 := cpabse.CphDec(s, pm.P)
	fmt.Println(cph1)

	t, _ := cpabse.CP_TkEnc(prv, keyword, msk, pm)
	fmt.Println("------------------------------------------")
	fmt.Println(t)
	fmt.Println("------------------------------------------")

	t1 := "58 255 129 3 1 1 6 66 121 116 101 84 107 1 255 130 0 1 4 1 4 84 111 107 49 1 10 0 1 4 84 111 107 50 1 10 0 1 4 84 111 107 51 1 10 0 1 5 67 111 109 112 115 1 255 134 0 0 0 35 255 133 2 1 1 20 91 93 42 99 112 97 98 115 101 46 66 121 116 101 84 107 67 111 109 112 1 255 134 0 1 255 132 0 0 39 255 131 3 1 2 255 132 0 1 3 1 4 65 116 116 114 1 12 0 1 4 65 95 106 49 1 10 0 1 4 66 95 106 49 1 10 0 0 0 254 4 178 255 130 1 255 128 45 244 136 166 192 108 132 247 47 141 255 167 215 156 215 142 49 151 152 125 240 226 4 196 110 182 113 74 95 109 85 153 148 101 82 123 76 181 43 146 155 236 126 10 30 122 34 250 70 84 168 189 194 216 28 12 71 9 235 242 24 153 160 224 162 240 235 87 23 186 61 114 89 73 225 23 183 10 193 107 25 244 175 12 219 91 26 31 205 103 197 62 10 24 213 247 241 162 58 247 179 213 105 198 66 17 111 155 138 152 56 72 245 155 134 105 121 159 140 10 255 93 51 214 113 198 193 24 1 255 128 78 5 112 122 133 220 14 168 140 111 76 47 118 177 152 241 49 83 84 160 72 163 124 57 209 42 119 183 192 244 168 67 63 4 146 224 82 220 149 91 8 234 139 233 1 144 191 157 170 137 109 200 25 26 11 20 136 220 58 137 193 70 105 126 45 217 37 46 96 73 177 251 255 208 22 12 108 233 118 145 252 163 230 106 50 117 128 161 128 91 90 110 214 226 190 155 16 56 12 151 36 255 130 55 131 27 157 192 237 250 102 55 189 90 144 204 151 252 74 81 9 199 157 143 187 79 45 64 1 255 128 4 54 181 108 91 49 102 120 74 231 14 159 64 138 36 133 13 32 232 232 255 29 56 130 228 4 150 213 25 37 7 251 171 50 54 217 133 31 221 218 239 78 0 204 51 2 186 81 38 52 60 245 97 145 7 106 206 84 81 193 26 70 128 223 49 249 137 223 224 182 233 22 26 217 213 16 56 127 70 76 0 255 175 225 178 51 23 180 171 20 5 23 63 150 144 218 184 135 79 175 167 33 229 216 94 78 149 142 122 203 64 143 78 52 185 184 191 85 48 164 115 215 119 190 192 133 190 249 1 3 1 3 102 111 111 1 255 128 126 87 76 178 240 196 18 215 83 145 249 102 184 165 188 199 194 38 27 52 138 237 219 183 156 235 74 71 195 80 7 52 128 20 90 86 75 7 211 123 158 82 104 68 8 4 156 195 186 209 176 226 167 212 215 119 186 212 202 217 12 39 60 210 39 128 126 109 40 80 217 196 7 35 166 128 66 231 40 8 165 230 101 244 187 140 201 159 138 82 124 175 83 132 39 4 21 178 200 151 107 66 94 24 46 46 46 103 163 199 69 138 98 179 174 31 160 251 119 235 5 44 143 60 126 150 128 164 1 255 128 129 14 169 141 180 29 166 124 227 117 27 224 180 86 181 38 163 81 195 83 37 132 194 88 186 68 102 210 82 198 121 211 43 55 106 5 131 4 138 172 231 182 25 231 16 245 15 101 59 44 118 201 162 43 79 147 117 178 62 182 151 97 43 47 66 150 87 226 180 120 240 153 65 83 199 162 22 109 48 52 220 194 170 82 149 242 133 197 55 198 161 6 153 104 158 237 46 171 0 174 108 130 55 61 233 170 148 96 86 111 157 62 120 163 192 66 52 140 127 121 130 253 16 216 208 153 246 163 0 1 3 102 105 109 1 255 128 136 223 93 112 198 201 153 165 152 140 8 23 69 228 41 19 89 147 189 148 199 113 151 119 78 248 0 147 215 43 53 203 138 123 143 134 124 222 145 195 170 90 215 18 106 152 179 161 253 81 87 37 25 236 98 60 14 159 144 155 84 119 230 159 11 151 87 154 213 83 36 222 97 187 49 151 198 195 163 206 130 93 180 153 120 112 250 173 84 169 56 241 97 196 114 146 127 129 93 5 83 175 216 128 186 182 128 93 217 166 255 158 133 94 79 102 38 77 68 23 209 250 133 175 202 70 253 139 1 255 128 100 105 243 247 205 232 137 113 63 172 246 20 228 165 6 251 56 7 20 131 80 122 6 154 157 137 211 41 148 210 100 9 55 18 152 62 46 237 95 19 162 95 129 21 198 95 81 160 11 29 195 230 211 202 249 227 107 198 81 116 179 61 213 103 130 206 160 162 245 119 22 214 238 79 93 237 173 229 179 55 29 255 52 98 153 155 130 222 36 124 206 117 188 77 249 252 224 29 87 249 3 143 249 248 165 179 160 65 158 232 19 19 119 78 50 214 36 146 46 218 212 91 161 38 179 195 179 115 0 1 3 98 97 102 1 255 128 26 111 139 78 30 155 35 161 228 63 57 10 186 46 121 61 190 254 136 248 48 161 135 96 123 171 154 69 54 167 101 123 215 26 44 65 99 207 216 250 212 3 144 8 31 124 123 0 92 231 251 244 113 89 186 17 166 167 229 196 124 11 81 70 4 56 249 179 248 86 7 70 251 87 134 94 16 218 28 62 62 93 128 209 84 242 29 77 196 193 24 159 9 240 235 213 65 209 150 208 186 78 143 47 176 96 16 40 3 196 236 21 162 201 242 139 40 201 143 48 87 35 181 176 115 112 214 27 1 255 128 103 249 133 131 188 141 56 255 39 170 182 5 145 184 96 116 249 186 207 1 238 182 124 181 95 122 145 161 6 108 116 163 124 131 216 215 22 205 86 23 173 7 110 157 250 103 63 31 111 69 203 126 102 57 65 11 203 26 215 224 211 198 124 139 52 170 143 251 200 87 48 154 201 203 28 188 51 122 203 116 86 244 78 198 65 233 72 42 26 113 183 41 171 32 129 188 252 232 112 27 156 84 91 199 24 54 8 100 218 203 250 221 96 12 132 240 62 223 58 108 195 110 53 52 212 190 103 206 0 0"

	tk1 := cpabse.TkDec(t1, pm.P)
	//fmt.Println(tk1)
	//c, cph := cpabse.CP_Enc(pm, policy, msk, keyword)
	//	fmt.Println(cph)
	//	cph1 := cpabse.CphDec(c, pm.P)
	//	fmt.Println(cph1)

	//_, tk := cpabse.CP_TkEnc(prv, keyword, msk, pm)

	if cpabse.Check(tk1, cph1, pm.P) {
		fmt.Println("ok")
	} else {
		fmt.Println("wrong")
	}

}